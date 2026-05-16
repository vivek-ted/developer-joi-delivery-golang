pipeline {
    agent any

    environment {
        IMAGE_NAME = "joi-delivery"
        REGISTRY   = "localhost:5000"   // local registry
        NAMESPACE  = "dev"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Docker build & push') {
            steps {
                sh """
                    docker build -t ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} .
                    docker push ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}
                """
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withKubeConfig([credentialsId: 'kubeconfig-local']) {
                    sh """
                        kubectl set image deployment/${IMAGE_NAME} \
                            ${IMAGE_NAME}=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} \
                            -n ${NAMESPACE}
                        kubectl rollout status deployment/${IMAGE_NAME} -n ${NAMESPACE}
                    """
                }
            }
        }
    }

    post {
        failure {
            withKubeConfig([credentialsId: 'kubeconfig-local']) {
                sh "kubectl rollout undo deployment/${IMAGE_NAME} -n ${NAMESPACE}"
            }
        }
    }
}