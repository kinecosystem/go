pipeline {
    agent any
    //parameters for the build
    parameters {
        //parameters for the environment creation

        //parameters for the load test
        string(name: 'VERSION', defaultValue: '1.0.0', description: 'tag/version for docjerhub image')
        string(name: 'BRANCH', defaultValue: 'jenkins', description: 'git branch (default: master)')
    }

    stages {
        stage ('Preparation'){
            steps {
                echo "Preparation"
                sh '''
                    rm -rf *
                    docker system prune -f
                '''
            }
        }
        stage ('Checkout code'){
            steps {
                echo "Checking out ${BRANCH}"
                // Install the desired Go version
                // def root = tool name: 'Go 1.11', type: 'go'

                // Export environment variables pointing to the directory where Go was installed
                    sh '''
                        mkdir -p go/src/github.com/kinecosystem/
                        cd go/src/github.com/kinecosystem
                        git clone -b ${BRANCH} https://github.com/kinecosystem/go.git
                    '''

                }
        }
        stage('Building') {
                steps {
                    echo "importing dependencies and building project"
                    sh '''
                        cd go/src/github.com/kinecosystem/go
                        make build
                    '''
                }
        }

        stage ('Run tests'){
                steps {
                    echo 'Running tests'
                    sh '''
                        cd go/src/github.com/kinecosystem/go
                        make test | tee a test_results.txt
                    '''
                }
        }
        stage ('Collecting and publishing results'){
                steps {
                    echo 'Collecting results'
                }
        }
        stage ('Packaging/Tagging'){
                steps {
                    echo 'packaging with version ${VERSION}'
                }
        }
        stage ('Pushing to Docker hub'){
                steps {
                    echo 'Pushing to Docker hub'
                }
        }
        stage ('Deploy to env'){
            steps{
                echo "Deploying to env -TBD"
            }
        }
    }
    post {
        always {
            echo 'Cleanup environment'
            sh '''
                cd go/src/github.com/kinecosystem/go
                make test_teardown
            '''
            junit '**/*.xml'
        }
    }
}
