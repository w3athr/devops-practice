pipeline {
    agent {
        node {
            label 'worker1-jobs'
        }
    }
    tools {
        go 'Go 1.25.0'
    }
    options {
        timestamps()
        disableConcurrentBuilds()
        gitLabConnection('yadro_gitlab_connection')
        gitlabBuilds(builds: ['lint', 'test', 'build', 'deploy'])
    }
    parameters {
        string(
            name: 'MANUAL_BRANCH',
            defaultValue: 'main',
            description: 'Ветка для ручного запуска. Для main оставь main.'
        )
    }   
    environment {
        AUTHOR = 'egor.volkov'
        VERSION = '0.5.0'
        SERVICE = 'weather'
        PORT = '8000'
        IMAGE_NAME = 'w3athr/weather-app'
        IMAGE_TAG = "build-${env.BUILD_NUMBER}"
    } 
    stages {
        stage('lint') {
            steps {
                gitlabCommitStatus('lint') {
                    sh '''
                    set -e
                    go version
                    go vet ./...
                    '''
                }
            }
        }
        stage('test') {
            steps {
                gitlabCommitStatus('test') {
                    sh '''
                    set -e
                    go test -v ./...
                    '''
                }
            }
        }
        stage('build') {
            steps {
                gitlabCommitStatus('build') {
                    withCredentials([usernamePassword(
                        credentialsId: 'dockerhub-pat',
                        usernameVariable: 'DOCKERHUB_USER',
                        passwordVariable: 'DOCKERHUB_PASS'
                    )]) {
                        sh '''
                        set -e
                        echo "$DOCKERHUB_PASS" | docker login -u "$DOCKERHUB_USER" --password-stdin
                        docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
                        docker push ${IMAGE_NAME}:${IMAGE_TAG}
                        docker logout
                        '''
                    }
                }
            }
        }
        stage('deploy') {
            when {
                expression {
                    return (env.gitlabBranch == 'main') || (!env.gitlabBranch && params.MANUAL_BRANCH == 'main')
                }
            }
        input {
            message "Deploy image ${IMAGE_NAME}:${IMAGE_TAG} to production?"
            ok "Deploy"
        }            
            environment {
                WEATHER_API_KEY = credentials('weather-api-key')
            }
            steps {
                gitlabCommitStatus('deploy') {
                    sh '''
                    set -e
                    printf '%s' "$WEATHER_API_KEY" > api_key.txt                    
                    chmod 600 api_key.txt
                    docker pull ${IMAGE_NAME}:${IMAGE_TAG}
                    docker compose -f docker-compose.hardened.yml up -d --force-recreate
                    curl -fsS http://127.0.0.1:${PORT}/info
                    '''
                }
            }
        }
    }
    post {
        always {
            sh 'rm -f api_key.txt || true'
        }
        success {
            echo 'Pipeline finished successfully'
        }
        failure {
            echo 'Pipeline failed'
        }
    }
}