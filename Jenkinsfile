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

        stage('Promotion & Deployment') {
            steps {
                script {
                    // 1. Build & Push Image
                    conditionalStage(
                        name: 'Build & Push Image',
                        condition: (env.BRANCH_NAME == 'main' || env.TAG_NAME || env.CHANGE_ID)
                    ) {
                        gitlabCommitStatus('build') {
                            def imageTag = env.TAG_NAME ?: "build-${env.BUILD_NUMBER}"
                            
                            withCredentials([usernamePassword(credentialsId: 'dockerhub_pat', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                                sh "docker login -u ${USER} -p ${PASS}"
                                sh "docker build -t w3athr/weather-app:${imageTag} ."
                                sh "docker push w3athr/weather-app:${imageTag}"
                            }
                        }
                    }

                    // 2. Deploy to Staging
                    conditionalStage(
                        name: 'Deploy to Staging',
                        condition: (env.BRANCH_NAME == 'main')
                    ) {
                        gitlabCommitStatus('deploy') {
                            build job: 'parameterized_pipeline', 
                                parameters: [
                                    string(name: 'IMAGE_TAG', value: "build-${env.BUILD_NUMBER}"),
                                    string(name: 'ENVIRONMENT', value: 'staging')
                                ]
                        }
                    }

                    // 3. Deploy to Production
                    conditionalStage(
                        name: 'Deploy to Production',
                        condition: (env.TAG_NAME != null)
                    ) {
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
    }

    post {
        failure {
            updateGitlabCommitStatus name: 'quality', state: 'failed'
        }
    }
}