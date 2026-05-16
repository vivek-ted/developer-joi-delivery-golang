pipeline {
    agent any

    environment {
        IMAGE_NAME = 'joi-delivery'
        REGISTRY   = 'localhost:5000'
        IMAGE_TAG  = "${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}"
        LATEST_TAG = "${REGISTRY}/${IMAGE_NAME}:latest"
        NAMESPACE  = 'dev'
    }

    stages {

        stage('Checkout') {
            steps {
                echo "Branch: ${env.BRANCH_NAME ?: 'unknown'}"
                checkout scm
            }
        }

        stage('Test') {
            steps {
                sh '''
                    export PATH=$PATH:/usr/local/go/bin
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

        stage('Docker build') {
            steps {
                sh '''
                    docker build \
                        -t ${IMAGE_TAG} \
                        -t ${LATEST_TAG} \
                        .
                '''
            }
        }

        stage('Push to local registry') {
            steps {
                sh '''
                    docker push ${IMAGE_TAG}
                    docker push ${LATEST_TAG}
                    echo "✅ Pushed ${IMAGE_TAG}"
                    echo "✅ Pushed ${LATEST_TAG}"
                '''
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withKubeConfig([credentialsId: 'kubeconfig-local']) {
                    sh '''
                        # ensure namespace exists
                        kubectl get namespace ${NAMESPACE} 2>/dev/null || \
                            kubectl create namespace ${NAMESPACE}

                        # apply all manifests
                        kubectl apply -f k8s/ -n ${NAMESPACE}

                        # roll out the new image
                        kubectl set image deployment/${IMAGE_NAME} \
                            ${IMAGE_NAME}=${IMAGE_TAG} \
                            -n ${NAMESPACE}

                        # wait up to 2 minutes for rollout
                        kubectl rollout status deployment/${IMAGE_NAME} \
                            -n ${NAMESPACE} --timeout=120s
                    '''
                }
            }
        }

        stage('Verify') {
            steps {
                withKubeConfig([credentialsId: 'kubeconfig-local']) {
                    sh '''
                        echo "──── Pods ────"
                        kubectl get pods -n ${NAMESPACE} -l app=${IMAGE_NAME}

                        echo "──── Service ────"
                        kubectl get svc -n ${NAMESPACE}

                        echo "──── Image in use ────"
                        kubectl get deployment ${IMAGE_NAME} -n ${NAMESPACE} \
                            -o jsonpath="{.spec.template.spec.containers[0].image}"
                        echo ""
                    '''
                }
            }
        }
    }

    post {
        success {
            echo "✅ Build #${BUILD_NUMBER} deployed to namespace '${NAMESPACE}'"
        }
        failure {
            echo "❌ Build #${BUILD_NUMBER} failed — rolling back"
            withKubeConfig([credentialsId: 'kubeconfig-local']) {
                sh 'kubectl rollout undo deployment/${IMAGE_NAME} -n ${NAMESPACE} || true'
            }
        }
        always {
            // clean up local image to save disk space
            sh 'docker rmi ${IMAGE_TAG} ${LATEST_TAG} || true'
            cleanWs()
        }
    }
}