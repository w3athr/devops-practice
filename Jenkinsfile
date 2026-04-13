@Library('jenkins_shared_library') _

pipeline {
    agent any

    options {
        timestamps()
        gitLabConnection('yadro_gitlab_connection') 
        gitlabBuilds(builds: ['quality', 'build', 'deploy'])
    }

    tools {
        go 'Go 1.25.0'
    }

    stages {
        stage('Static Checks') {
            steps {
                script {
                    gitlabCommitStatus('quality') {
                        parallel(
                            "Linting": {
                                sh 'go vet ./...'
                            },
                            "SAST Scan": {
                                runSAST.runSAST()
                            }
                        )
                    }
                }
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
                    branch 'main'
                    tag 'v*'
                    branch 'e.volkov/upgrade_Jenkinsfile'
                    expression { env.CHANGE_ID != null } 
                }
            }
            steps {
                script {
                    gitlabCommitStatus('build') {
                        def imageTag = env.TAG_NAME ?: "build-${env.BUILD_NUMBER}"
                        
                        withCredentials([usernamePassword(credentialsId: 'dockerhub_pat', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                            sh "docker login -u ${USER} -p ${PASS}"
                            sh "docker build -t w3athr/weather-app:${imageTag} ."
                            sh "docker push w3athr/weather-app:${imageTag}"
                        }
                    }
                }
            }
        }

        stage('Deploy to Staging') {
            when {
                anyOf {
                    branch 'main'
                    branch 'e.volkov/upgrade_Jenkinsfile'
                }
            }
            steps {
                script {
                    gitlabCommitStatus('deploy') {
                        build job: 'deploy-job', 
                            parameters: [
                                string(name: 'IMAGE_TAG', value: "build-${env.BUILD_NUMBER}"),
                                string(name: 'ENVIRONMENT', value: 'staging')
                            ]
                    }
                }
            }
        }

        stage('Deploy to Production') {
            when { tag "v*" } 
            steps {
                script {
                    gitlabCommitStatus('deploy') {
                        build job: 'parameterized_pipeline', 
                            parameters: [
                                string(name: 'IMAGE_TAG', value: "${env.TAG_NAME}"),
                                string(name: 'ENVIRONMENT', value: 'production')
                            ]
                    }
                }
            }
        }
    }

    post {
        failure {
            updateGitlabCommitStatus name: 'quality', state: 'failed'
        }
    }
}