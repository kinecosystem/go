pipeline {
    agent any
    //parameters for the build
    parameters {
        //parameters for the environment creation

        //parameters for the load test
        string(name: 'BRANCH', defaultValue: 'jenkins', description: 'git branch (default: master)')
        string(name: 'VERSION', defaultValue: '', description: 'docker and horizon version \n Overrides automatic version')

    }

    stages {
        stage ('Cleanup'){
            steps {
                echo "Environment cleanup"
                sh '''
                    rm -rf *
                    docker system prune -f
                    docker rmi horizon_horizon || true
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
        stage ('Run tests'){
            steps {
                echo 'Running tests and coverage'
                sh '''
                    cd go/src/github.com/kinecosystem/go
                    make test
                '''
            }
        }
        stage('Building horizon') {
            steps {
                echo "building horizon executable with version ${VERSION}"
                sh '''
                    VERSION="${VERSION}" \
                    TARGET="builder"
                    make build
                    '''
                echo 'copying horizon executable to local filesystem'
                sh '''
                    make get_horizon
                '''
            }
        }
        stage ('Copying horizon'){
            steps {
                echo 'copying horizon executable to local filesystem'
            }
        }
        stage ('Tagging and Pushing to Docker hub'){
            steps {
                echo 'Pushing Docker image (version: ${VERSION}) to dockerhub'
                echo "Creating release image with version ${VERSION}"
                sh '''
                    VERSION="${VERSION}" \
                    make docker_release
                '''
                withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh '''
                        VERSION="${VERSION}" \
                        make docker_push
                    '''
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
            echo "publishing reports"
            junit '**/test-results.xml'
            echo 'Cleanup environment'
            sh '''
                cd go/src/github.com/kinecosystem/go
                MOUNT_POINT=${MOUNT_POINT} make tests_teardown
                MOUNT_POINT=${MOUNT_POINT} make jenkins_teardown
            '''

        }
    }
}
