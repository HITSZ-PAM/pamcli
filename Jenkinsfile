pipeline {
	agent any
	stages {
		stage('Check Go Version') {
			def root = tool type: 'go', name: 'Go Build Env'

    		// Export environment variables pointing to the directory where Go was installed
    		withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
      		  	sh 'go version'
    		}
		}
		stage('Build') {
			def root = tool type: 'go', name: 'Go Build Env'

    		// Export environment variables pointing to the directory where Go was installed
    		withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
      		  	sh 'go build'
    		}
		}
		stage('Deploy') {
			steps {
				sh 'echo $DEMO_PASSWORD'
			}
		}
	}
}