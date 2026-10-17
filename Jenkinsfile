pipeline {
    agent any

    options {
        timeout(time: 60, unit: 'MINUTES')
        ansiColor('xterm')
        disableConcurrentBuilds()
    }

    parameters {
        choice(
            name: 'TARGET_SERVICE',
            choices: [
                'all',
                'admission-service',
                'auth-service',
                'bed-service',
                'billing-service',
                'document-service',
                'emr-service',
                'laboratory-service',
                'master-service',
                'patient-service',
                'pharmacy-service',
                'practitioner-service',
                'radiology-service',
                'bff-web',
                'bff-mobile',
                'b2b-gateway',
                'satusehat-worker',
                'bpjs-worker',
                'notification-worker',
                'migration-runner'
            ],
            description: 'Target microservice to test and build (select "all" for complete monorepo)'
        )
        booleanParam(name: 'RUN_PROTO_LINT', defaultValue: true, description: 'Run Buf Protobuf linter')
        booleanParam(name: 'RUN_LINT', defaultValue: true, description: 'Run golangci-lint on target service(s)')
        booleanParam(name: 'RUN_TESTS', defaultValue: true, description: 'Execute unit and integration test suite')
        booleanParam(name: 'BUILD_DOCKER', defaultValue: true, description: 'Build Docker container image(s)')
        string(name: 'DOCKER_REGISTRY', defaultValue: 'selaras', description: 'Docker image namespace / registry prefix')
        string(name: 'IMAGE_TAG', defaultValue: '', description: 'Custom image tag (defaults to BUILD_NUMBER if empty)')
    }

    environment {
        TAG_NAME = "${params.IMAGE_TAG ?: env.BUILD_NUMBER}"
        REGISTRY = "${params.DOCKER_REGISTRY}"
    }

    stages {
        stage('Verify Environment') {
            steps {
                script {
                    echo "=========================================================="
                    echo "Selaras SIMRS CI/CD Pipeline"
                    echo "Target Service : ${params.TARGET_SERVICE}"
                    echo "Build Number   : ${env.BUILD_NUMBER}"
                    echo "Image Tag      : ${env.TAG_NAME}"
                    echo "=========================================================="
                    runCmd('go version')
                    runCmd('docker --version')
                }
            }
        }

        stage('Protobuf Validation') {
            when {
                expression { return params.RUN_PROTO_LINT }
            }
            steps {
                script {
                    echo "--> Running Buf Protobuf Lint..."
                    runCmd('buf lint')
                }
            }
        }

        stage('Code Quality (golangci-lint)') {
            when {
                expression { return params.RUN_LINT }
            }
            steps {
                script {
                    echo "--> Running golangci-lint..."
                    def servicesToLint = getServiceList(params.TARGET_SERVICE)
                    for (svc in servicesToLint) {
                        echo "Linting ${svc}..."
                        if (isUnix()) {
                            sh "(cd ${svc} && golangci-lint run ./...)"
                        } else {
                            bat "cd ${svc} && golangci-lint run ./..."
                        }
                    }
                }
            }
        }

        stage('Execute Tests') {
            when {
                expression { return params.RUN_TESTS }
            }
            steps {
                script {
                    echo "--> Executing Test Suite..."
                    runCmd('go work sync')
                    def servicesToTest = getServiceList(params.TARGET_SERVICE)
                    for (svc in servicesToTest) {
                        echo "Testing ${svc}..."
                        runCmd("go test -v -cover ./${svc}/...")
                    }
                }
            }
        }

        stage('Verify Compilation') {
            steps {
                script {
                    echo "--> Verifying Go binary compilation..."
                    def servicesToBuild = getServiceList(params.TARGET_SERVICE)
                    for (svc in servicesToBuild) {
                        if (svc == 'scripts/migrator' || svc == 'migration-runner') {
                            runCmd("go build ./scripts/migrator")
                        } else {
                            runCmd("go build ./${svc}/cmd")
                        }
                    }
                }
            }
        }

        stage('Build Docker Images') {
            when {
                expression { return params.BUILD_DOCKER }
            }
            steps {
                script {
                    echo "--> Building Docker Images..."
                    def servicesToDocker = getDockerServiceList(params.TARGET_SERVICE)
                    for (item in servicesToDocker) {
                        def img = "${env.REGISTRY}/${item.name}:${env.TAG_NAME}"
                        def latestImg = "${env.REGISTRY}/${item.name}:latest"
                        echo "Building ${img} (Context: ${item.context}, Dockerfile: ${item.dockerfile})..."

                        def buildCmd = "docker build -f ${item.dockerfile} -t ${img} -t ${latestImg} ${item.context}"
                        runCmd(buildCmd)
                    }
                }
            }
        }
    }

    post {
        always {
            echo "Pipeline finished with status: ${currentBuild.currentResult}"
        }
        success {
            echo "SUCCESS: Selaras SIMRS pipeline for target '${params.TARGET_SERVICE}' passed successfully!"
        }
        failure {
            echo "FAILURE: Selaras SIMRS pipeline failed. Please inspect build logs."
        }
    }
}

// Helper function to execute shell or batch command based on agent OS
def runCmd(String cmd) {
    if (isUnix()) {
        sh cmd
    } else {
        bat cmd
    }
}

// Return array of service paths for Go commands
def getServiceList(String target) {
    def allServices = [
        'admission-service',
        'auth-service',
        'bed-service',
        'billing-service',
        'document-service',
        'emr-service',
        'laboratory-service',
        'master-service',
        'patient-service',
        'pharmacy-service',
        'practitioner-service',
        'radiology-service',
        'bff-web',
        'bff-mobile',
        'b2b-gateway',
        'satusehat-worker',
        'bpjs-worker',
        'notification-worker',
        'scripts/migrator'
    ]

    if (target == 'all') {
        return allServices
    } else if (target == 'migration-runner') {
        return ['scripts/migrator']
    } else {
        return [target]
    }
}

// Return Docker build metadata (name, context, and dockerfile path)
def getDockerServiceList(String target) {
    def catalog = [
        [name: 'admission-service',   context: './admission-service', dockerfile: './admission-service/Dockerfile'],
        [name: 'auth-service',        context: './auth-service',      dockerfile: './auth-service/Dockerfile'],
        [name: 'bed-service',         context: './bed-service',       dockerfile: './bed-service/Dockerfile'],
        [name: 'billing-service',     context: './billing-service',   dockerfile: './billing-service/Dockerfile'],
        [name: 'document-service',    context: './document-service',  dockerfile: './document-service/Dockerfile'],
        [name: 'emr-service',         context: './emr-service',       dockerfile: './emr-service/Dockerfile'],
        [name: 'laboratory-service',  context: './laboratory-service',dockerfile: './laboratory-service/Dockerfile'],
        [name: 'master-service',      context: './master-service',    dockerfile: './master-service/Dockerfile'],
        [name: 'patient-service',     context: './patient-service',   dockerfile: './patient-service/Dockerfile'],
        [name: 'pharmacy-service',    context: './pharmacy-service',  dockerfile: './pharmacy-service/Dockerfile'],
        [name: 'practitioner-service',context: './practitioner-service',dockerfile: './practitioner-service/Dockerfile'],
        [name: 'radiology-service',   context: './radiology-service', dockerfile: './radiology-service/Dockerfile'],
        [name: 'bff-web',             context: '.',                   dockerfile: 'bff-web/Dockerfile'],
        [name: 'bff-mobile',          context: '.',                   dockerfile: 'bff-mobile/Dockerfile'],
        [name: 'b2b-gateway',         context: '.',                   dockerfile: 'b2b-gateway/Dockerfile'],
        [name: 'satusehat-worker',    context: '.',                   dockerfile: 'satusehat-worker/Dockerfile'],
        [name: 'bpjs-worker',         context: '.',                   dockerfile: 'bpjs-worker/Dockerfile'],
        [name: 'notification-worker', context: '.',                   dockerfile: 'notification-worker/Dockerfile'],
        [name: 'migration-runner',    context: '.',                   dockerfile: './scripts/migrator/Dockerfile']
    ]

    if (target == 'all') {
        return catalog
    }

    for (item in catalog) {
        if (item.name == target || (target == 'migration-runner' && item.name == 'migration-runner')) {
            return [item]
        }
    }

    return []
}