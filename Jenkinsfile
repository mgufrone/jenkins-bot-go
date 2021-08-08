@Library('shared-pipeline')_
import com.mgufron.SlackSectionedMessage
import com.mgufron.RunInLog

slack = new SlackSectionedMessage(this, "#general")
runner = new RunInLog(slack, true)
slack.sendMessage(slack.message("Building started: <${env.BUILD_URL}|${env.BUILD_TAG}>"))
try {
  runner.run(message: "Setup") {
    podTemplate(inheritFrom: "golang sonar") {
      node(POD_LABEL) {
        checkout(scm: scm, changelog: true).each { k, v ->
          env.setProperty(k, v)
        }
        stage("Build") {
          try {
            container("golang") {
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
  }
  runner.run(message: "Deployment") {
    podTemplate(inheritFrom: "golang helm builder") {
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
  }
  slack.success()
} catch (e) {
  slack.error()
}