@Library('jenkins_shared_library') _

node {
    properties([gitLabConnection('yadro_gitlab_connection')])

    timestamps { 
        def goHome = tool 'Go 1.25.0'
        
        withEnv(["PATH+GO=${goHome}/bin"]) {
            try {
                stage('Checkout') {
                    checkout scm
                }

                stage('Static Checks') {
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

                def isTag = (env.TAG_NAME != null)
                def isMain = (env.BRANCH_NAME == 'main')
                def isMR = (env.CHANGE_ID != null)
                
                conditionalStage(
                    name: 'Build & Push Image',
                    condition: (isMain || isTag || isMR),
                    gitlabStatus: 'build'
                ) {
                    def imageTag = env.TAG_NAME ?: (isMR ? "mr-${env.CHANGE_ID}" : "build-${env.BUILD_NUMBER}")
                    
                    withCredentials([usernamePassword(credentialsId: 'dockerhub_pat', usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                        sh "docker login -u ${USER} -p ${PASS}"
                        sh "docker build -t w3athr/weather-app:${imageTag} ."
                        sh "docker push w3athr/weather-app:${imageTag}"
                    }
                }

                conditionalStage(
                    name: isTag ? "Deploy to Production" : "Deploy to Staging",
                    condition: (isTag || isMain),
                    gitlabStatus: 'deploy'
                ) {
                    def targetEnv = isTag ? 'production' : 'staging'
                    def targetTag = isTag ? env.TAG_NAME : "build-${env.BUILD_NUMBER}"

                    build job: 'parameterized_pipeline',
                        parameters: [
                            string(name: 'IMAGE_TAG', value: targetTag),
                            string(name: 'ENVIRONMENT', value: targetEnv)
                        ]
                }

            } catch (Exception e) {
                updateGitlabCommitStatus name: 'quality', state: 'failed'
                echo "Pipeline failed with error: ${e.message}"
                throw e 
            } finally {
                archiveArtifacts artifacts: 'sast-report.json', allowEmptyArchive: true
            }
        }
    }
}