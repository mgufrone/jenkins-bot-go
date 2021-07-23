pipeline {
  agent {
    kubernetes {
      inheritFrom "deployment sonar"
      yaml '''
spec:
    containers:
    - name: golang
      image: golang:1.16
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
        script {
          try {
            copyArtifacts(projectName: 'jenkins-bot', selectors: lastSuccessful);
            sh "ls -lah"
          } catch (Exception e) {
            echo "something wrong: $e"
          }
        }
        container("golang") {
          sh "go get gotest.tools/gotestsum"
          sh "go mod vendor"
          sh "go vet ./..."
          sh "gotestsum --format dots --junitfile report.xml"
        }
      }
      post {
        success {
          archiveArtifacts artifacts: "vendor/**/*", fingerprint: true
        }
        always {
          junit "report.xml"
          container('sonar') {
            sh('sonar-scanner -Dsonar.login=$SONAR_LOGIN')
          }
        }
      }
    }
    stage("Deployment") {
      steps {
        container('kaniko') {
          script {
            env.IMAGE_TAG = env.GIT_COMMIT.subString(0, 6)
            def data = ["auths": ["ghcr.io": ["username": env.GITHUB_USR, "password": env.GITHUB_PSW]]]
            writeJSON file: "docker-config.json", json: data
            sh "cp docker-config.json /kaniko/.docker/config.json"
            sh "/kaniko/executor --context . --dockerfile ./Dockerfile --destination ghcr.io/mgufrone/jenkins-bot:${env.GIT_BRANCH} --destination ghcr.io/mgufrone/jenkins-bot:${env.IMAGE_TAG}"
          }
        }
        container('helm') {
          sh "helm upgrade --install jenkins-bot ./jenkin-bots --set image.tag=${env.IMAGE_TAG}"
        }
      }
    }
  }
}