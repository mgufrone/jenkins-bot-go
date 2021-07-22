pipeline {
  agent {
    kubernetes {
      inheritFrom "deployment"
      yaml '''
spec:
    containers:
    - name: golang
      image: golang:1.16-alpine
      command:
      - sleep
      args:
      - infinity
'''
    }
  }
  environment {
    SONAR_HOST_URL = credentials('sonar-url')
    SONAR_LOGIN = credentials('sonar-token')
    GITHUB = credentials('github')
  }
  stages {
    stage("Build") {
      steps {
        container("golang") {
          sh "go mod vendor"
          sh "go vet ./..."
          sh "go test ./... -o coverage.out"
        }
      }
    }
    stage("Deployment") {
      steps {
        container('kaniko') {
          script {
            def data = ["auths": ["ghcr.io": ["username": env.GITHUB_USR, "password": env.GITHUB_PSW]]]
            writeJSON file: "docker-config.json", json: data
            sh "cp docker-config.json /kaniko/.docker/config.json"
            sh "cp composer.production.json composer.json"
            sh "cp composer.production.lock composer.lock"
            sh "/kaniko/executor --context . --dockerfile ./build.Dockerfile --destination ghcr.io/mgufrone/jenkins-bot:${env.GIT_BRANCH} --destination ghcr.io/mgufrone/symfony-test:${env.GIT_COMMIT}"
          }
        }
        container('helm') {
          sh "helm upgrade --install jenkins-bot ./jenkin-bots --set image.tag=${env.GIT_COMMIT}"
        }
      }
    }
  }
}