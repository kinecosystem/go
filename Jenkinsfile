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
                def root = tool name: 'Go 1.11', type: 'go'

                // Export environment variables pointing to the directory where Go was installed
                withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
                    sh '''
                        go version
                        mkdir -p $GOPATH/src/github.com/kinecosystem
                        cd $_
                        git clone -b ${BRANCH} https://github.com/kinecosystem/go.git
                    '''

                }
            }
        }
        stage('Import dependencies') {
                steps {
                    echo "importing dependencies"
                    sh '''
                        export GOPATH=$PWD/go
                        cd $GOPATH/src/github.com/kinecosystem/go
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
