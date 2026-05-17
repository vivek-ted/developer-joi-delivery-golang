pipeline {
    agent any

    environment {
        IMAGE_NAME  = 'joi-delivery'
        REGISTRY    = 'localhost:5000'
        IMAGE_TAG   = "${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}"
        LATEST_TAG  = "${REGISTRY}/${IMAGE_NAME}:latest"
        APP_PORT    = '8001'
        HOST_PORT   = '9090'           // access app at http://localhost:9090
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
                echo "✅ Checked out: ${env.GIT_COMMIT?.take(7) ?: 'unknown'}"
            }
        }

        stage('Test') {
            steps {
                sh '''
                    export PATH=$PATH:/usr/local/go/bin
                    export GOPATH=/home/jenkins/go
                    export GOCACHE=/home/jenkins/.cache/go-build
                    go version
                    go mod tidy
                    go test ./... -v
                '''
            }
        }

        stage('Build image') {
            steps {
                sh """
                    docker build -t ${IMAGE_TAG} -t ${LATEST_TAG} .
                    echo "✅ Image built: ${IMAGE_TAG}"
                """
            }
        }

        stage('Push to registry') {
            steps {
                sh """
                    docker push ${IMAGE_TAG}
                    docker push ${LATEST_TAG}
                    echo "✅ Pushed to local registry"
                """
            }
        }

        stage('Deploy') {
            steps {
                sh """
                    # Stop and remove existing container if running
                    docker stop ${IMAGE_NAME} 2>/dev/null || true
                    docker rm   ${IMAGE_NAME} 2>/dev/null || true

                    # Pull latest image and run
                    docker pull ${LATEST_TAG}
                    docker run -d \
                        --name ${IMAGE_NAME} \
                        --restart unless-stopped \
                        -p ${HOST_PORT}:${APP_PORT} \
                        ${LATEST_TAG}

                    echo "✅ Container started"
                    docker ps | grep ${IMAGE_NAME}
                """
            }
        }

        stage('Health check') {
            steps {
                sh """
                    # Wait for app to start
                    sleep 3

                    # Check container is running
                    docker inspect -f '{{.State.Running}}' ${IMAGE_NAME}

                    echo "✅ App running at http://localhost:${HOST_PORT}"
                    echo "✅ Health: http://localhost:${HOST_PORT}/health"
                """
            }
        }
    }

    post {
        success {
            echo "✅ Build #${BUILD_NUMBER} deployed — http://localhost:${HOST_PORT}"
        }
        failure {
            echo "❌ Build #${BUILD_NUMBER} failed — check console output above"
            sh "docker stop ${IMAGE_NAME} 2>/dev/null || true"
        }
        always {
            sh "docker rmi ${IMAGE_TAG} 2>/dev/null || true"
            deleteDir()
        }
    }
}