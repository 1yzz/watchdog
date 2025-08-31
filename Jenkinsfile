pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21'
        NODE_VERSION = '18'
        IMAGE_NAME = 'watchdog'
    }
    
    tools {
        go "${GO_VERSION}"
        nodejs "${NODE_VERSION}"
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
        
        stage('Install Dependencies') {
            parallel {
                stage('Go Dependencies') {
                    steps {
                        sh 'go mod download'
                    }
                }
                
                stage('Node Dependencies') {
                    steps {
                        dir('sdk/javascript') {
                            sh 'npm install'
                        }
                    }
                }
            }
        }
        
        stage('Build') {
            parallel {
                stage('Build Server') {
                    steps {
                        sh '''
                            make proto-generate
                            make ent-generate
                            make build
                        '''
                        archiveArtifacts artifacts: 'bin/watchdog'
                    }
                }
                
                stage('Build SDK') {
                    steps {
                        dir('sdk/javascript') {
                            sh 'npm run build'
                        }
                    }
                }
            }
        }
        
        stage('Test') {
            parallel {
                stage('Go Tests') {
                    steps {
                        sh 'go test ./...'
                    }
                }
                
                stage('JS Tests') {
                    steps {
                        dir('sdk/javascript') {
                            sh 'npm test'
                        }
                    }
                }
            }
        }
        
        stage('Docker Build') {
            steps {
                script {
                    docker.build("${IMAGE_NAME}:${BUILD_TAG}")
                }
            }
        }
        
        stage('Deploy Local') {
            when {
                branch 'main'
            }
            steps {
                script {
                    echo 'Deploying to local Docker...'
                    
                    // Stop existing container if running
                    sh '''
                        docker stop watchdog-local || true
                        docker rm watchdog-local || true
                    '''
                    
                    // Run new container
                    sh """
                        docker run -d \\
                            --name watchdog-local \\
                            --restart unless-stopped \\
                            -p 50051:50051 \\
                            -p 8080:8080 \\
                            -v /var/lib/watchdog:/data \\
                            -e DATABASE_URL=sqlite:///data/watchdog.db \\
                            -e LOG_LEVEL=info \\
                            ${IMAGE_NAME}:${BUILD_TAG}
                    """
                    
                    // Wait for service to be ready
                    sh '''
                        echo "Waiting for service to start..."
                        sleep 10
                        
                        # Health check
                        if curl -f http://localhost:8080/health; then
                            echo "✅ Watchdog deployed successfully!"
                        else
                            echo "❌ Deployment failed - service not responding"
                            exit 1
                        fi
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
