Feature: As a MLOps engineer I would like the Model Registry to store metadata information about models
  Taken from User Stories

  Scenario: As a MLOps engineer I would like to store Model name
    Given I have a connection to MR
    When I store a VersionedModel with payload:
      """
      {"model_name": "Pricing Model", "version":"v1", "model_uri":"s3://12345"}
      """
    Then there should be a mlmd Context of type "odh.RegisteredModel" named "Pricing Model"
    And there should be a mlmd Context of type "odh.ModelVersion" having property named "model_name" valorised with string value "Pricing Model"

