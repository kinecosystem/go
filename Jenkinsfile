pipeline {
    agent any
    //parameters for the build
    parameters {
        //parameters for the environment creation

        //parameters for the load test
        string(name: 'BRANCH', defaultValue: 'jenkins', description: 'git branch (default: master)')
        string(name: 'VERSION', defaultValue: '', description: 'docker and horizon version \n Overrides automatic versioning')
        string(name: 'MOUNT_POINT', defaultValue: '"/jenkins_home/workspace/horizon/go/src/github.com/kinecosystem/go/"', \
            description: 'local mount point for docker images')

    }

    stages {
        stage ('Cleanup'){
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
                sh '''
                    mkdir -p go/src/github.com/kinecosystem/
                    cd go/src/github.com/kinecosystem
                    git clone -b ${BRANCH} https://github.com/kinecosystem/go.git
                '''
                }
        }
        stage ('Unit tests'){
            steps{
                echo "placeholder for unit tests"
            }
        }
        stage('Building docker image') {
            steps {
                echo "importing dependencies and building docker image with version ${VERSION}"
                // DATE='1-1-1970' VERSION="1.2.3" MOUNT_POINT="/Home/ubuntu/go/src/github.com/kinecosystem/go/" make build
                //TODO: get version automatically
                //TODO: get vendor (kinecosystem) as parameter
                sh '''
                    DATE=`date +%F-%T`
                    VERSION="0.0.1-dev"
                    DATE="${DATE}" VERSION="${VERSION}" MOUNT_POINT="${MOUNT_POINT}" make build
                '''
            }
        }
        stage ('Run tests'){
            steps {
                echo 'Running tests and coverage'
                sh '''
                    cd go/src/github.com/kinecosystem/go
                    MOUNT_POINT=${MOUNT_POINT} make test
                '''
            }
        }
        stage ('Collecting and publishing results'){
            steps {
                echo 'Collecting results - TBD'
            }
        }
        stage ('Pushing to Docker hub'){
            steps {
                echo 'Pushing Docker image, version $env.GIT_REVISION to dockerhub'
                withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh 'make push'
                }
            }
        }
        stage ('Deploy to env'){
            steps{
                echo "Deploying to env - TBD"
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
            junit 'test-results.xml'
        }
    }
}
