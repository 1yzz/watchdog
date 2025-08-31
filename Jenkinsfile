pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21'
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
                    # Check if the correct Go version is installed
                    CURRENT_GO_VERSION=$(go version 2>/dev/null | grep -o 'go[0-9.]*' | sed 's/go//' || echo "none")
                    if [ "$CURRENT_GO_VERSION" != "${GO_VERSION}" ]; then
                        echo "Installing Go ${GO_VERSION}... (current: $CURRENT_GO_VERSION)"
                        wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                        sudo rm -rf /usr/local/go
                        sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                        rm go${GO_VERSION}.linux-amd64.tar.gz
                    else
                        echo "Go ${GO_VERSION} already installed"
                    fi
                    
                    # Set Go environment
                    export PATH=$PATH:/usr/local/go/bin
                    export GOPATH=$HOME/go
                    export PATH=$PATH:$GOPATH/bin
                    
                    # Verify Go installation
                    go version
                    
                    # Ensure we're using the exact version
                    INSTALLED_VERSION=$(go version | grep -o 'go[0-9.]*' | sed 's/go//')
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
                    go mod download
                    
                    # Install protoc-gen-go and protoc-gen-go-grpc
                    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
                    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
                    
                    # Verify tools
                    protoc --version
                    go version
                '''
            }
        }
        
        stage('Build Server') {
            steps {
                sh '''
                    # Generate protobuf code
                    make proto-generate
                    
                    # Generate Ent database code  
                    make ent-generate
                    
                    # Build the Go server binary
                    make build
                    
                    # Verify binary
                    ls -la bin/watchdog
                    file bin/watchdog
                '''
                archiveArtifacts artifacts: 'bin/watchdog'
            }
        }
        
        stage('Test') {
            steps {
                sh '''
                    # Run Go tests
                    go test -v ./...
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
                    echo "🚀 Deploying Watchdog locally..."
                    
                    # Stop existing process if running
                    pkill -f "watchdog" || true
                    
                    # Create data directory
                    mkdir -p /var/lib/watchdog
                    
                    # Copy binary to deployment location
                    sudo cp bin/watchdog /usr/local/bin/watchdog
                    sudo chmod +x /usr/local/bin/watchdog
                    
                    # Start the service in background
                    nohup /usr/local/bin/watchdog \\
                        --grpc-port=50051 \\
                        --http-port=8080 \\
                        --database-url="sqlite:///var/lib/watchdog/watchdog.db" \\
                        --log-level=info > /var/log/watchdog.log 2>&1 &
                    
                    echo $! > /var/run/watchdog.pid
                    
                    # Wait for service to start
                    echo "Waiting for service to start..."
                    sleep 10
                    
                    # Health check with retry
                    echo "Performing health check..."
                    for i in {1..12}; do
                        if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
                            echo "✅ Watchdog deployed successfully!"
                            echo "🔗 gRPC endpoint: localhost:50051"
                            echo "🔗 HTTP endpoint: http://localhost:8080"
                            echo "📝 Logs: /var/log/watchdog.log"
                            echo "📊 PID: $(cat /var/run/watchdog.pid)"
                            exit 0
                        fi
                        echo "Health check attempt ${i}/12 failed, retrying in 5s..."
                        sleep 5
                    done
                    
                    echo "❌ Deployment failed - service not responding after 60s"
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
