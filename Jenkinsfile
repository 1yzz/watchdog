pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.24.6'
        CGO_ENABLED = '0'
        GOOS = 'linux'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.BUILD_TAG = "${env.BUILD_NUMBER}-${sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()}"
                }
            }
        }
        
        stage('Setup Build Environment') {
            steps {
                sh '''
                    # Check if Go is already installed
                    if command -v go${GO_VERSION} &> /dev/null; then
                        echo "Go ${GO_VERSION} already installed, checking version..."
                        export PATH=/usr/local/go${GO_VERSION}/bin:$PATH
                        INSTALLED_VERSION=$(go${GO_VERSION} version | grep -o 'go[0-9.]*' | sed 's/go//' | tr -d '\n')
                        
                        if [ "$INSTALLED_VERSION" = "${GO_VERSION}" ]; then
                            echo "✅ Go ${GO_VERSION} already installed and correct version"
                        else
                            echo "❌ Go ${GO_VERSION} installed but wrong version: ${INSTALLED_VERSION}"
                            exit 1
                        fi
                    else
                        # Install exact Go version
                        echo "Installing Go ${GO_VERSION}..."
                        wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                        sudo rm -rf /usr/local/go${GO_VERSION}
                        sudo tar -C /usr/local --transform="s/^go/go${GO_VERSION}/" -xzf go${GO_VERSION}.linux-amd64.tar.gz
                        rm go${GO_VERSION}.linux-amd64.tar.gz
                        
                        # Create version-specific symlink
                        sudo ln -sf /usr/local/go${GO_VERSION}/bin/go /usr/local/bin/go${GO_VERSION}
                        
                        # Set Go environment
                        export PATH=/usr/local/go${GO_VERSION}/bin:$PATH
                        export GOPATH=$HOME/go
                        export PATH=$PATH:$GOPATH/bin
                        
                        # Verify Go installation
                        go${GO_VERSION} version
                        
                        # Ensure we're using the exact version
                        INSTALLED_VERSION=$(go${GO_VERSION} version | grep -o 'go[0-9.]*' | sed 's/go//' | tr -d '\n')
                        if [ "$INSTALLED_VERSION" != "${GO_VERSION}" ]; then
                            echo "❌ Expected Go ${GO_VERSION}, but got $INSTALLED_VERSION"
                            exit 1
                        fi
                        echo "✅ Go ${GO_VERSION} installed successfully"
                    fi
                    
                    # Install protoc if not available
                    if ! command -v protoc &> /dev/null; then
                        echo "Installing protoc..."
                        wget -q https://github.com/protocolbuffers/protobuf/releases/download/v24.4/protoc-24.4-linux-x86_64.zip
                        sudo unzip -q protoc-24.4-linux-x86_64.zip -d /usr/local
                        rm protoc-24.4-linux-x86_64.zip
                    fi
                    
                    # Download Go dependencies
                    go${GO_VERSION} mod download
                    
                    # Install protoc-gen-go and protoc-gen-go-grpc
                    go${GO_VERSION} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
                    go${GO_VERSION} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
                    
                    # Verify tools
                    protoc --version
                    go${GO_VERSION} version
                '''
            }
        }
        
        stage('Build Server') {
            steps {
                sh '''
                    # Set Go environment
                    export PATH=/usr/local/go${GO_VERSION}/bin:$PATH
                    export GOPATH=$HOME/go
                    export PATH=$PATH:$GOPATH/bin
                    
                    # Verify Go tools are available
                    which go${GO_VERSION}
                    which protoc-gen-go || echo "protoc-gen-go not found, installing..."
                    which protoc-gen-go-grpc || echo "protoc-gen-go-grpc not found, installing..."
                    
                    # Install Go protobuf tools if not available
                    if ! command -v protoc-gen-go &> /dev/null; then
                        go${GO_VERSION} install google.golang.org/protobuf/cmd/protoc-gen-go@latest
                    fi
                    
                    if ! command -v protoc-gen-go-grpc &> /dev/null; then
                        go${GO_VERSION} install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
                    fi
                    
                    # Verify tools are now available
                    protoc-gen-go --version
                    protoc-gen-go-grpc --version
                    
                    # Generate protobuf code
                    make proto
                    
                    # Generate Ent database code  
                    make ent-generate
                    
                    # Build the Go server binary
                    make build
                    
                    # Verify binary
                    ls -la bin/watchdog-server
                    file bin/watchdog-server
                '''
                archiveArtifacts artifacts: 'bin/watchdog-server'
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                    # Set Go environment
                    export PATH=/usr/local/go${GO_VERSION}/bin:$PATH
                    
                    # Run Go tests
                    go${GO_VERSION} test -v ./...
                '''
            }
        }
        
        stage('Deploy Local') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                }
            }
            steps {
                sh '''
                    # Set Go environment for grpcurl
                    export PATH=/usr/local/go${GO_VERSION}/bin:$PATH
                    export GOPATH=$HOME/go
                    export PATH=$PATH:$GOPATH/bin
                    
                    echo "🚀 Deploying Watchdog locally using systemd..."
                    
                    # Stop existing service if running
                    sudo systemctl stop watchdog.service || true
                    sudo systemctl disable watchdog.service || true
                    
                    # Kill any process using port 50051
                    echo "Checking for processes using port 50051..."
                    PORT_PID=$(sudo lsof -ti:50051 2>/dev/null || sudo netstat -tlnp 2>/dev/null | grep :50051 | awk '{print $7}' | cut -d'/' -f1 || echo "")
                    if [ ! -z "$PORT_PID" ]; then
                        echo "Found process $PORT_PID using port 50051, killing it..."
                        sudo kill -9 $PORT_PID || true
                        sleep 2
                    fi
                    
                    # Also try to kill any watchdog processes that might be running
                    echo "Killing any existing watchdog processes..."
                    sudo pkill -f "watchdog-server" || true
                    sleep 1
                    
                    # Create data directory with proper permissions
                    sudo mkdir -p /var/lib/watchdog
                    sudo chown jenkins:jenkins /var/lib/watchdog
                    sudo chmod 755 /var/lib/watchdog
                    
                    # Copy binary to data directory
                    sudo cp bin/watchdog-server /var/lib/watchdog/watchdog-server
                    sudo chmod +x /var/lib/watchdog/watchdog-server
                    
                    # Copy systemd service file if it doesn't exist
                    if [ ! -f /etc/systemd/system/watchdog.service ]; then
                        echo "Installing watchdog systemd service..."
                        sudo cp scripts/watchdog.service /etc/systemd/system/watchdog.service
                        sudo chmod 644 /etc/systemd/system/watchdog.service
                        sudo systemctl daemon-reload
                        sudo systemctl enable watchdog.service
                    else
                        echo "Watchdog systemd service already exists, skipping installation"
                    fi
                    
                    # Restart service
                    sudo systemctl restart watchdog.service
                    
                    # Wait for service to start
                    echo "Waiting for service to start..."
                    sleep 10
                    
                    # Check service status
                    echo "Service status:"
                    sudo systemctl status watchdog.service --no-pager -l
                    
                    # gRPC health check with retry
                    echo "Performing gRPC health check..."
                    for i in {1..12}; do
                        if grpcurl -plaintext 127.0.0.1:50051 watchdog.WatchdogService/GetHealth > /dev/null 2>&1; then
                            echo "✅ Watchdog deployed successfully!"
                            echo "🔗 gRPC endpoint: localhost:50051"
                            echo "📝 Logs: journalctl -u watchdog.service"
                            echo "📊 Service: watchdog.service"
                            
                            # Show health status
                            echo "🏥 Health status:"
                            grpcurl -plaintext 127.0.0.1:50051 watchdog.WatchdogService/GetHealth
                            exit 0
                        fi
                        echo "gRPC health check attempt ${i}/12 failed, retrying in 5s..."
                        sleep 5
                    done
                    
                    echo "❌ Deployment failed - gRPC service not responding after 60s"
                    echo "Service logs:"
                    sudo journalctl -u watchdog.service --no-pager -n 20
                    exit 1
                '''
            }
            post {
                failure {
                    sh '''
                        echo "🔍 Deployment failed, collecting debug info..."
                        sudo systemctl status watchdog.service --no-pager -l || true
                        sudo journalctl -u watchdog.service --no-pager -n 50 || true
                    '''
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        
        success {
            echo 'Build succeeded!'
        }
        
        failure {
            echo 'Build failed!'
        }
    }
}
