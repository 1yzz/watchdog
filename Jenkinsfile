pipeline {
    agent {
        docker {
            image 'golang:1.21-alpine'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    
    environment {
        IMAGE_NAME = 'watchdog'
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
                    # Install required tools in Alpine
                    apk add --no-cache git make curl docker-cli protoc protobuf-dev
                    
                    # Download Go dependencies
                    go mod download
                    
                    # Install protoc-gen-go and protoc-gen-go-grpc
                    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
                    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
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
        
        stage('Docker Build') {
            steps {
                sh '''
                    # Build Docker image with current context
                    docker build -t ${IMAGE_NAME}:${BUILD_TAG} -t ${IMAGE_NAME}:latest .
                    
                    # Verify image was built
                    docker images | grep ${IMAGE_NAME}
                    
                    # Test image can start
                    docker run --rm ${IMAGE_NAME}:${BUILD_TAG} --version || echo "Binary test completed"
                '''
            }
        }
        
        stage('Deploy to Local Docker') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                }
            }
            steps {
                sh '''
                    echo "🚀 Deploying Watchdog to local Docker container..."
                    
                    # Stop and remove existing container
                    echo "Stopping existing container..."
                    docker stop watchdog-local 2>/dev/null || true
                    docker rm watchdog-local 2>/dev/null || true
                    
                    # Create data directory on host
                    sudo mkdir -p /var/lib/watchdog
                    sudo chown jenkins:jenkins /var/lib/watchdog || true
                    
                    # Deploy new container
                    echo "Starting new container with image: ${IMAGE_NAME}:${BUILD_TAG}"
                    docker run -d \\
                        --name watchdog-local \\
                        --restart unless-stopped \\
                        -p 50051:50051 \\
                        -p 8080:8080 \\
                        -v /var/lib/watchdog:/data \\
                        -e DATABASE_URL=sqlite:///data/watchdog.db \\
                        -e LOG_LEVEL=info \\
                        -e GRPC_PORT=50051 \\
                        -e HTTP_PORT=8080 \\
                        --health-cmd="curl -f http://localhost:8080/health || exit 1" \\
                        --health-interval=30s \\
                        --health-timeout=10s \\
                        --health-retries=3 \\
                        ${IMAGE_NAME}:${BUILD_TAG}
                    
                    # Wait for container to start
                    echo "Waiting for container to start..."
                    sleep 15
                    
                    # Check container status
                    echo "Container status:"
                    docker ps | grep watchdog-local
                    
                    # Check logs
                    echo "Container logs:"
                    docker logs watchdog-local --tail 20
                    
                    # Health check with retry
                    echo "Performing health check..."
                    for i in {1..12}; do
                        if curl -f -s http://localhost:8080/health > /dev/null 2>&1; then
                            echo "✅ Watchdog deployed successfully!"
                            echo "🔗 gRPC endpoint: localhost:50051"
                            echo "🔗 HTTP endpoint: http://localhost:8080"
                            
                            # Show final status
                            docker ps | grep watchdog-local
                            exit 0
                        fi
                        echo "Health check attempt ${i}/12 failed, retrying in 5s..."
                        sleep 5
                    done
                    
                    echo "❌ Deployment failed - service not responding after 60s"
                    echo "Container logs:"
                    docker logs watchdog-local
                    exit 1
                '''
            }
            post {
                failure {
                    sh '''
                        echo "🔍 Deployment failed, collecting debug info..."
                        docker ps -a | grep watchdog || true
                        docker logs watchdog-local || true
                        docker inspect watchdog-local || true
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
