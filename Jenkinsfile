podTemplate(inheritFrom: "golang sonar") {
  node(POD_LABEL) {
  checkout(scm: scm, changelog: true).each { k, v ->
    env.setProperty(k, v)
  }
  stage("Build") {
    try {
      container("golang") {
        sh "go get gotest.tools/gotestsum"
        sh "go mod vendor"
        sh "gotestsum --format testname --junitfile report.xml"
      }
      withSonarQubeEnv("sonar") {
        container('sonar') {
          sh('sonar-scanner -Dsonar.qualitygate.wait=true -Dsonar.qualitygate.timeout=600')
        }
      }
      stash name: env.BUILD_TAG, includes: "**", useDefaultExcludes: false
    } finally {
      junit "report.xml"
    }
  }
  }
}
podTemplate(inheritFrom: "builder golang helm") {
node(POD_LABEL) {
  stage("Deployment") {
    unstash "${env.BUILD_TAG}"
  }
  container("golang") {
    sh "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ."
  }
  container('kaniko') {
    script {
      env.IMAGE_TAG = env.GIT_COMMIT.substring(0, 6)
      sh "/kaniko/executor --context . --dockerfile ./Dockerfile --destination ghcr.io/mgufrone/jenkins-bot:${env.GIT_BRANCH} --destination ghcr.io/mgufrone/jenkins-bot:${env.IMAGE_TAG}"
    }
  }
  container('helm') {
    sh "helm upgrade --install jenkins-bot ./jenkin-bots --set image.tag=${env.IMAGE_TAG}"
  }
  }
}