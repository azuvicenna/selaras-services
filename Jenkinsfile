pipeline {
    agent any

    environment {
        SERVICE = 'patient-service'
        IMAGE   = "selaras/${SERVICE}"
    }

    stages {
        stage('Test') {
            steps {
                bat "go test ./${SERVICE}/..."
            }
        }

        stage('Build Image') {
            steps {
                bat "docker build -t ${IMAGE}:${BUILD_NUMBER} ${SERVICE}"
            }
        }
    }
}