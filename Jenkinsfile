pipeline {
    agent any

    parameters {
        choice(name: 'ENVIRONMENT', choices: ['staging', 'production'], description: 'Select environment to deploy', defaultValue: 'staging')
    }

    environment {
        // Define environment variables based on the selected environment
        GO_VERSION = '1.21.3'
        DOCKER_COMPOSE_VERSION = '2.21.0'
        DOCKER_COMPOSE_FILE = "${params.ENVIRONMENT == 'production' ? 'docker-compose.prod.yaml' : 'docker-compose.staging.yaml'}"
    }

    stages {
        stage('Checkout') {
            steps {
                // This stage checks out the source code from your version control system
                checkout scm
            }
        }

        stage('Run Docker Compose') {
            steps {
                // Run Docker Compose to start your application and any required services
                script {
                    def dockerComposeCmd = "docker compose up -d"
                    sh dockerComposeCmd
                    echo "INFO: Successfully deployed to ${params.ENVIRONMENT} server"
                }
            }
        }
    }

    post {
        success {
            // This block is executed if the pipeline is successful
            echo 'Pipeline succeeded!'
        }

        failure {
            // This block is executed if the pipeline fails
            echo 'Pipeline failed!'
        }
    }
}
