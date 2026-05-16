pipeline {
    agent any

    environment {
        IMAGE_NAME    = 'joi-delivery'
        REGISTRY      = 'your-dockerhub-username'
        IMAGE_TAG     = "${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}"
        NAMESPACE     = 'dev'
        KUBECONFIG_ID = 'kubeconfig-local'
    }

    stages {

        stage('Checkout') {
            steps {
                echo "Checking out branch: ${env.BRANCH_NAME}"
                checkout scm
            }
        }

        stage('Test') {
            steps {
                sh '''
                    go version
                    go mod tidy
                    go test ./... -v -coverprofile=coverage.out
                    go tool cover -func=coverage.out
                '''
            }
            post {
                always {
                    archiveArtifacts artifacts: 'coverage.out', allowEmptyArchive: true
                }
            }
        }

        stage('Build Docker image') {
            steps {
                sh "docker build -t ${IMAGE_TAG} ."
            }
        }

        stage('Push Docker image') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'dockerhub-creds',
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )
                ]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker push ${IMAGE_TAG}
                        docker logout
                    '''
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withKubeConfig([credentialsId: "${KUBECONFIG_ID}"]) {
                    sh '''
                        # Ensure namespace exists
                        kubectl get namespace ${NAMESPACE} || \
                            kubectl create namespace ${NAMESPACE}

                        # Apply manifests
                        kubectl apply -f k8s/ -n ${NAMESPACE}

                        # Update image to this build
                        kubectl set image deployment/${IMAGE_NAME} \
                            ${IMAGE_NAME}=${IMAGE_TAG} \
                            -n ${NAMESPACE}

                        # Wait for rollout to complete
                        kubectl rollout status deployment/${IMAGE_NAME} \
                            -n ${NAMESPACE} --timeout=120s
                    '''
                }
            }
        }

        stage('Verify deployment') {
            steps {
                withKubeConfig([credentialsId: "${KUBECONFIG_ID}"]) {
                    sh '''
                        echo "--- Pods ---"
                        kubectl get pods -n ${NAMESPACE} -l app=${IMAGE_NAME}

                        echo "--- Service ---"
                        kubectl get svc -n ${NAMESPACE}

                        echo "--- Recent events ---"
                        kubectl get events -n ${NAMESPACE} \
                            --sort-by=.metadata.creationTimestamp | tail -10
                    '''
                }
            }
        }
    }

    post {
        success {
            echo "✅ Deployment successful — build ${BUILD_NUMBER}"
        }

        failure {
            echo "❌ Build failed — rolling back..."
            withKubeConfig([credentialsId: "${KUBECONFIG_ID}"]) {
                sh '''
                    kubectl rollout undo deployment/${IMAGE_NAME} -n ${NAMESPACE}
                    kubectl rollout status deployment/${IMAGE_NAME} -n ${NAMESPACE}
                '''
            }
        }

        always {
            sh 'docker rmi ${IMAGE_TAG} || true'
            cleanWs()
        }
    }
}