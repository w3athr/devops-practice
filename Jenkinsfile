@Library('jenkins_shared_library') _

pipeline {
    agent any

    options {
        timestamps()
        gitLabConnection('yadro_gitlab_connection') 
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
                                runSAST()
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
                        condition: (env.BRANCH_NAME == 'main' || env.TAG_NAME || env.CHANGE_ID),
                        gitlabStatus: 'build' 
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
                    // 2. Deploy to Prod/Staging
                    def isTag = (env.TAG_NAME != null)
                    def isMain = (env.BRANCH_NAME == 'main')

                    conditionalStage(
                        name: isTag ? "Deploy to Production" : "Deploy to Staging",
                        condition: (isTag || isMain),
                        gitlabStatus: 'deploy'
                    ) {
                        def targetEnv = isTag ? 'production' : 'staging'
                        def targetTag = isTag ? env.TAG_NAME : "build-${env.BUILD_NUMBER}"

                        gitlabCommitStatus('deploy') {
                            build job: 'parameterized_pipeline',
                                parameters: [
                                    string(name: 'IMAGE_TAG', value: targetTag),
                                    string(name: 'ENVIRONMENT', value: targetEnv)
                                ]
                        }
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