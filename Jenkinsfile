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
                    condition: (isTag || isMain || env.BRANCH_NAME == 'e.volkov/argocd-delivery-add'),
                    gitlabStatus: 'deploy'
                ) {
                    def targetEnv = isTag ? 'production' : 'staging'
                    def targetTag = isTag ? env.TAG_NAME : "build-${env.BUILD_NUMBER}"
                    def gitopsRepo = "education-git.yadro.com/education/devops/2026/e.volkov/helm_for_argocd.git"

                    echo "Updating GitOps repository for ${targetEnv}..."

                    withCredentials([usernamePassword(credentialsId: 'gitlab-gitops-token', usernameVariable: 'GIT_USER', passwordVariable: 'GIT_PASS')]) {
                        sh """
                            # Чистим старые клоны и клонируем репо с чартом
                            rm -rf helm_for_argocd
                            git clone https://${GIT_USER}:${GIT_PASS}@${gitopsRepo}
                            cd helm_for_argocd/weather-app

                            # Меняем тег образа в соответствующем файле values
                            # Ищем строку 'tag: "..."' и заменяем на новый тег
                            sed -i 's/tag: .*/tag: "${targetTag}"/' values-${targetEnv}.yaml

                            # Фиксируем изменения в Git
                            git config user.email "jenkins-bot@yadro.com"
                            git config user.name "Jenkins GitOps Bot"
                            git add values-${targetEnv}.yaml
                                
                            # [skip ci] нужен, чтобы не зациклить сборку, если пайплайны в одном репо
                            git commit -m "chore(cd): update ${targetEnv} image to ${targetTag} [skip ci]"
                            git push origin main
                            """

                    echo "Successfully updated GitOps repo. ArgoCD will sync shortly."
                    }
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