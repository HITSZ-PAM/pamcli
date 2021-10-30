pipeline {
	agent any
	tools {
		go 'Go Build Env'
	}
	stages {
		stage('Check Go Version') {
    		steps {
      		  	sh 'go version'
    		}
		}
		stage('Build') {
    		steps {
      		  	sh 'go build'
    		}
		}
		stage('Deploy') {
			environment {
				DEMO_PASSWORD = credentials('demo_password')
			}
			steps {
				sh 'echo $DEMO_PASSWORD'
			}
		}
	}
}