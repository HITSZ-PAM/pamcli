node {
	stage('Check Go Version') {
		def root = tool type: 'go', name: 'Go Build Env'
   		// Export environment variables pointing to the directory where Go was installed
   		steps {
   		  	sh 'go version'
   		}
	}
	stage('Build') {
		def root = tool type: 'go', name: 'Go Build Env'
    		// Export environment variables pointing to the directory where Go was installed
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