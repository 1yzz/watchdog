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
                    echo "✅ Using Go ${GO_VERSION}"
                    
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
                    
                    echo "🚀 Deploying Watchdog locally..."
                    
                    # Stop existing process if running
                    pkill -f "watchdog-server" || true
                    
                    # Create data directory with proper permissions
                    sudo mkdir -p /var/lib/watchdog
                    sudo chown jenkins:jenkins /var/lib/watchdog
                    
                    # Create log directory with proper permissions
                    sudo mkdir -p /var/log
                    sudo touch /var/log/watchdog.log
                    sudo chown jenkins:jenkins /var/log/watchdog.log
                    
                    # Copy binary to deployment location
                    sudo cp bin/watchdog-server /usr/local/bin/watchdog-server
                    sudo chmod +x /usr/local/bin/watchdog-server
                    
                    # Start the service in background (gRPC only)
                    nohup /usr/local/bin/watchdog-server > /var/log/watchdog.log 2>&1 &
                    
                    echo $! > /tmp/watchdog.pid
                    
                    # Wait for service to start
                    echo "Waiting for service to start..."
                    sleep 10
                    
                    # gRPC health check with retry
                    echo "Performing gRPC health check..."
                    for i in {1..12}; do
                        if grpcurl -plaintext 127.0.0.1:50051 watchdog.WatchdogService/GetHealth > /dev/null 2>&1; then
                            echo "✅ Watchdog deployed successfully!"
                            echo "🔗 gRPC endpoint: localhost:50051"
                            echo "📝 Logs: /var/log/watchdog.log"
                            echo "📊 PID: $(cat /tmp/watchdog.pid)"
                            
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
                    tail -20 /var/log/watchdog.log
                    exit 1
                '''
            }
            post {
                failure {
                    sh '''
                        echo "🔍 Deployment failed, collecting debug info..."
                        ps aux | grep watchdog || true
                        tail -50 /var/log/watchdog.log || true
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
