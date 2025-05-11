package probe

const (
	ListProbeQuery = `query ListProbes($projectID: ID!, $probeNames: [ID!], $filter: ProbeFilterInput) {
		listProbes(projectID: $projectID, probeNames: $probeNames, filter: $filter) {
		  name
		  type
		  createdAt
		  createdBy{
			username
		  }
		}
	  }
	`
	GetProbeQuery = `query getProbe($projectID: ID!, $probeName: ID!) {
		getProbe(projectID: $projectID, probeName: $probeName) {
		  name
		  description
		  type
		  infrastructureType
		  kubernetesHTTPProperties{
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
			url
			insecureSkipVerify
			method {
				get {
					responseCode
					criteria
				}
				post {
					body
					contentType
					responseCode
					criteria
				}
			}
		  }
		  kubernetesCMDProperties{
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
			command
			comparator {
				type
				criteria
				value
			}
			source
		  }
		  k8sProperties {
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
			group
			version
			resource
			namespace
		  }
		  promProperties {
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
			endpoint
			query
		  }
		  createdAt
		  createdBy{
			username
		  }
		  updatedAt
		  updatedBy{
			username
		  }
		  tags
		}
	  }
	`
	GetProbeYAMLQuery = `query getProbeYAML($projectID: ID!, $request: GetProbeYAMLRequest!) {
		getProbeYAML(projectID: $projectID, request: $request)
	  }
	`

	DeleteProbeQuery = `mutation deleteProbe($probeName: ID!, $projectID: ID!) {
		deleteProbe(probeName: $probeName, projectID: $projectID)
	  }
	`
	createProbeMutation = `mutation addProbe($projectID: ID!, $request: ProbeRequest!) {
		addProbe(projectID: $projectID, request: $request) {
			name
			description
			type
			infrastructureType
			tags
			kubernetesHTTPProperties {
				probeTimeout
				interval
				url
				method {
					get {
						responseCode
						criteria
					}
					post {
						body
						contentType
						responseCode
						criteria
					}
				}
				insecureSkipVerify
			}
			kubernetesCMDProperties {
				probeTimeout
				interval
				command
				comparator {
					type
					criteria
					value
				}
			}
			k8sProperties {
				probeTimeout
				interval
				retry
				group
				version
				resource
				namespace
			}
			promProperties {
				probeTimeout
				interval
				endpoint
				query
			}
		}
	}`
)
