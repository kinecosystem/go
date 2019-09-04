pipeline {
    agent any
    //parameters for the build
    parameters {
        //parameters for the environment creation

        //parameters for the load test
        string(name: 'VERSION', defaultValue: '1.0.0', description: 'tag/version for docjerhub image')
        string(name: 'BRANCH', defaultValue: 'master', description: 'git branch (default: master)')
    }
    stages {
        stage ('Checkout code'){
            steps {
                echo "Checking out ${BRANCH}"
                sh '''
                    git clone -b ${BRANCH} https://github.com/kinecosystem/go.git
                '''
            }
        }
        stage('Import dependencies') {
                steps {
                    echo "importing dependencies"
                    sh '''
                        make dep
                    '''
                }
        }
        stage ('Building'){
                steps {
                    echo 'building project"}'
                }
        }

        stage ('Run tests'){
                steps {
                    echo 'Running tests'

                }
        }
        stage ('Collecting results'){
                steps {
                    echo 'Collecting results'
                }
        }
        stage ('Packaging'){
                steps {
                    echo 'packaging with version ${VERSION}'
                }
        }
        stage ('Pushing to Docker hub'){
                steps {
                    echo 'Pushing to Docker hub"}'
                }
        }
    }
}
