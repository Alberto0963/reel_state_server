package models

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var DB *gorm.DB

// func ConnectDataBase(){

// 	err := godotenv.Load(".env")

// 	if err != nil {
// 	  log.Fatalf("Error loading .env file")
// 	}

// 	Dbdriver := os.Getenv("DB_DRIVER")
// 	DbHost := os.Getenv("DB_HOST")
// 	DbUser := os.Getenv("DB_USER")
// 	DbPassword := os.Getenv("DB_PASSWORD")
// 	DbName := os.Getenv("DB_NAME")
// 	DbPort := os.Getenv("DB_PORT")

// 	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

// 	DB, err = gorm.Open(Dbdriver, DBURL)

// 	if err != nil {
// 		fmt.Println("Cannot connect to database ", Dbdriver)
// 		log.Fatal("connection error:", err)
// 	} else {
// 		fmt.Println("We are connected to the database ", Dbdriver)
// 	}

// 	DB.AutoMigrate(&User{})

// }

var Pool *gorm.DB // Declare a global variable to hold the connection pool

// SSH tunnel configuration
type SSHTunnelConfig struct {
	SSHHost     string
	SSHPort     string
	SSHUser     string
	SSHPassword string
	SSHKeyPath  string
	SSHKeyPassphrase string 
	LocalPort   string
	RemoteHost  string
	RemotePort  string
}

// createSSHTunnel establishes an SSH tunnel for database connection
func createSSHTunnel(config SSHTunnelConfig) (*ssh.Client, net.Listener, error) {
    // SSH client configuration
    sshConfig := &ssh.ClientConfig{
        User:            config.SSHUser,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, use proper host key verification
        Timeout:         30 * time.Second,
    }

    // Authentication method (password or key)
    if config.SSHPassword != "" {
        sshConfig.Auth = []ssh.AuthMethod{
            ssh.Password(config.SSHPassword),
        }
    } else if config.SSHKeyPath != "" {
        key, err := os.ReadFile(config.SSHKeyPath)
        if err != nil {
            return nil, nil, fmt.Errorf("unable to read private key: %v", err)
        }

        var signer ssh.Signer
        if config.SSHKeyPassphrase != "" {
            // Key with passphrase
            signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(config.SSHKeyPassphrase))
        } else {
            // Key without passphrase
            signer, err = ssh.ParsePrivateKey(key)
        }
        
        if err != nil {
            return nil, nil, fmt.Errorf("unable to parse private key: %v", err)
        }

        sshConfig.Auth = []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        }
    } else {
        return nil, nil, fmt.Errorf("either SSH password or key path must be provided")
    }

    // Connect to SSH server
    sshClient, err := ssh.Dial("tcp", net.JoinHostPort(config.SSHHost, config.SSHPort), sshConfig)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to SSH server: %v", err)
    }

    // Create local listener
    localListener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", config.LocalPort))
    if err != nil {
        sshClient.Close()
        return nil, nil, fmt.Errorf("failed to create local listener: %v", err)
    }

    // Start forwarding connections
    go func() {
        for {
            localConn, err := localListener.Accept()
            if err != nil {
                return
            }

            go func(conn net.Conn) {
                defer conn.Close()

                // Connect to remote MySQL server through SSH tunnel
                remoteConn, err := sshClient.Dial("tcp", net.JoinHostPort(config.RemoteHost, config.RemotePort))
                if err != nil {
                    fmt.Printf("Failed to connect to remote MySQL server: %v\n", err)
                    return
                }
                defer remoteConn.Close()

                // Forward traffic bidirectionally between local and remote connections
                done := make(chan bool, 2)

                // Copy from local to remote
                go func() {
                    defer func() { done <- true }()
                    _, _ = io.Copy(remoteConn, conn)
                }()

                // Copy from remote to local
                go func() {
                    defer func() { done <- true }()
                    _, _ = io.Copy(conn, remoteConn)
                }()

                // Wait for either direction to finish
                <-done
            }(localConn)
        }
    }()

    fmt.Printf("SSH tunnel established: localhost:%s -> %s:%s -> %s:%s\n",
        config.LocalPort, config.SSHHost, config.SSHPort, config.RemoteHost, config.RemotePort)

    return sshClient, localListener, nil
}

func InitDB() (*gorm.DB, error) {

	baseDir, err := os.Getwd() // Get the current working directory

	if err != nil {
		log.Fatalf("Error loading path")

	}

	err = godotenv.Load(baseDir + "/.env")

	if err != nil {
		// log.Fatalf("Error loading .env file: ", err)
		fmt.Println("Error loading .env file:", err)

	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	// SSH tunnel configuration (optional)
	useSSHTunnel := os.Getenv("USE_SSH_TUNNEL")
	var sshClient *ssh.Client
	var localListener net.Listener

    if useSSHTunnel == "true" {
        sshConfig := SSHTunnelConfig{
            SSHHost:     os.Getenv("SSH_HOST"),
            SSHPort:     os.Getenv("SSH_PORT"),
            SSHUser:     os.Getenv("SSH_USER"),
            SSHPassword: os.Getenv("SSH_PASSWORD"),
            SSHKeyPath:  os.Getenv("SSH_KEY_PATH"),
            SSHKeyPassphrase: os.Getenv("SSH_KEY_PASSPHRASE"), // Agregar esta línea
            LocalPort:   os.Getenv("SSH_LOCAL_PORT"),
            RemoteHost:  os.Getenv("SSH_REMOTE_HOST"),
            RemotePort:  os.Getenv("SSH_REMOTE_PORT"),
        }

		// Set default values if not provided
		if sshConfig.SSHPort == "" {
			sshConfig.SSHPort = "22"
		}
		if sshConfig.LocalPort == "" {
			sshConfig.LocalPort = "3307" // Use different port to avoid conflicts
		}
		if sshConfig.RemoteHost == "" {
			sshConfig.RemoteHost = "127.0.0.1"
		}
		if sshConfig.RemotePort == "" {
			sshConfig.RemotePort = "3306"
		}

		fmt.Println("Establishing SSH tunnel...")
		sshClient, localListener, err = createSSHTunnel(sshConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create SSH tunnel: %v", err)
		}

		// Update connection parameters to use local tunnel
		DbHost = "127.0.0.1"

		// Convert local port to string if it's a number
		if localPort, err := strconv.Atoi(sshConfig.LocalPort); err == nil {
			DbPort = strconv.Itoa(localPort)
		} else {
			DbPort = sshConfig.LocalPort
		}

		// Give the tunnel a moment to establish
		time.Sleep(2 * time.Second)
	}

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	// dsn := os.Getenv(DBURL)

	db, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{})
	if err != nil {
		// Clean up SSH tunnel if connection fails
		if sshClient != nil {
			sshClient.Close()
		}
		if localListener != nil {
			localListener.Close()
		}
		return nil, err
	} else {
		fmt.Println("We are connected to the database ", Dbdriver)
		if useSSHTunnel == "true" {
			fmt.Println("Connection established through SSH tunnel")
		}
	}
	db = db.Debug()


	sqlDB, err := db.DB()
	if err != nil {
		// Clean up SSH tunnel if getting DB fails
		if sshClient != nil {
			sshClient.Close()
		}
		if localListener != nil {
			localListener.Close()
		}
		return nil, err
	}

	// Set the maximum number of connections in the pool
	sqlDB.SetMaxOpenConns(10)
	Pool = db
	return db, nil
}
