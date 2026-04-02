pipeline {
    agent {
        node {
            label 'worker1-jobs'
        }
    }
    stages {
        stage('Build') {
            steps {
                echo "Building..."
                sh '''
                echo "doing build stuff"
                '''
            }
        }
    }
}