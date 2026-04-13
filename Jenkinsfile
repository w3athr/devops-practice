@Library('jenkins_shared_library') _

pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Static Checks') {
            steps {
                parallel(
                    "Linting": {
                        sh 'go vet ./...'
                    },
                    "SAST Scan": {
                        script {
                            common.runSAST()
                        }
                    }
                )
            }
            post {
                always {
                    archiveArtifacts artifacts: 'sast-report.json', allowEmptyArchive: true
                }
            }
        }

        stage('Build & Push Image') {
            when {
                anyOf {
                    branch 'main'; tag 'v*'
                    expression { env.CHANGE_ID != null } 
                }
            }
            steps {
                script {
                    def imageTag = env.TAG_NAME ?: "build-${env.BUILD_NUMBER}"
                    
                    withCredentials([usernamePassword(credentialsId: 'dockerhub_pat', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                        sh "docker login -u ${USER} -p ${PASS}"
                        sh "docker build -t w3athr/weather-app:${imageTag} ."
                        sh "docker push w3athr/weather-app:${imageTag}"
                    }
                }
            }
        }

        stage('Deploy to Staging') {
            when { branch 'main' }
            steps {
                build job: 'deploy-job', 
                    parameters: [
                        string(name: 'IMAGE_TAG', value: "build-${env.BUILD_NUMBER}"),
                        string(name: 'ENVIRONMENT', value: 'staging')
                    ]
            }
        }

        stage('Deploy to Production') {
            when { tag "v*" } 
            steps {
                build job: 'deploy-job', 
                    parameters: [
                        string(name: 'IMAGE_TAG', value: "${env.TAG_NAME}"),
                        string(name: 'ENVIRONMENT', value: 'production')
                    ]
            }
        }
    }
}