pipeline {
	agent any
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
	}
}