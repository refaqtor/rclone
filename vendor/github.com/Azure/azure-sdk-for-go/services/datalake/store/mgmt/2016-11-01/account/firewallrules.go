package account

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// FirewallRulesClient is the creates an Azure Data Lake Store account management client.
type FirewallRulesClient struct {
	BaseClient
}

// NewFirewallRulesClient creates an instance of the FirewallRulesClient client.
func NewFirewallRulesClient(subscriptionID string) FirewallRulesClient {
	return NewFirewallRulesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewFirewallRulesClientWithBaseURI creates an instance of the FirewallRulesClient client.
func NewFirewallRulesClientWithBaseURI(baseURI string, subscriptionID string) FirewallRulesClient {
	return FirewallRulesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate creates or updates the specified firewall rule. During update, the firewall rule with the specified
// name will be replaced with this new firewall rule.
//
// resourceGroupName is the name of the Azure resource group. accountName is the name of the Data Lake Store
// account. firewallRuleName is the name of the firewall rule to create or update. parameters is parameters
// supplied to create or update the firewall rule.
func (client FirewallRulesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, parameters CreateOrUpdateFirewallRuleParameters) (result FirewallRule, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.CreateOrUpdateFirewallRuleProperties", Name: validation.Null, Rule: true,
				Chain: []validation.Constraint{{Target: "parameters.CreateOrUpdateFirewallRuleProperties.StartIPAddress", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.CreateOrUpdateFirewallRuleProperties.EndIPAddress", Name: validation.Null, Rule: true, Chain: nil},
				}}}}}); err != nil {
		return result, validation.NewError("account.FirewallRulesClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, accountName, firewallRuleName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client FirewallRulesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, parameters CreateOrUpdateFirewallRuleParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"firewallRuleName":  autorest.Encode("path", firewallRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/firewallRules/{firewallRuleName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client FirewallRulesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client FirewallRulesClient) CreateOrUpdateResponder(resp *http.Response) (result FirewallRule, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes the specified firewall rule from the specified Data Lake Store account
//
// resourceGroupName is the name of the Azure resource group. accountName is the name of the Data Lake Store
// account. firewallRuleName is the name of the firewall rule to delete.
func (client FirewallRulesClient) Delete(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(ctx, resourceGroupName, accountName, firewallRuleName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client FirewallRulesClient) DeletePreparer(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"firewallRuleName":  autorest.Encode("path", firewallRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/firewallRules/{firewallRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client FirewallRulesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client FirewallRulesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the specified Data Lake Store firewall rule.
//
// resourceGroupName is the name of the Azure resource group. accountName is the name of the Data Lake Store
// account. firewallRuleName is the name of the firewall rule to retrieve.
func (client FirewallRulesClient) Get(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string) (result FirewallRule, err error) {
	req, err := client.GetPreparer(ctx, resourceGroupName, accountName, firewallRuleName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client FirewallRulesClient) GetPreparer(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"firewallRuleName":  autorest.Encode("path", firewallRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/firewallRules/{firewallRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client FirewallRulesClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client FirewallRulesClient) GetResponder(resp *http.Response) (result FirewallRule, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByAccount lists the Data Lake Store firewall rules within the specified Data Lake Store account.
//
// resourceGroupName is the name of the Azure resource group. accountName is the name of the Data Lake Store
// account.
func (client FirewallRulesClient) ListByAccount(ctx context.Context, resourceGroupName string, accountName string) (result FirewallRuleListResultPage, err error) {
	result.fn = client.listByAccountNextResults
	req, err := client.ListByAccountPreparer(ctx, resourceGroupName, accountName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "ListByAccount", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.frlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "ListByAccount", resp, "Failure sending request")
		return
	}

	result.frlr, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "ListByAccount", resp, "Failure responding to request")
	}

	return
}

// ListByAccountPreparer prepares the ListByAccount request.
func (client FirewallRulesClient) ListByAccountPreparer(ctx context.Context, resourceGroupName string, accountName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/firewallRules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByAccountSender sends the ListByAccount request. The method will close the
// http.Response Body if it receives an error.
func (client FirewallRulesClient) ListByAccountSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByAccountResponder handles the response to the ListByAccount request. The method always
// closes the http.Response Body.
func (client FirewallRulesClient) ListByAccountResponder(resp *http.Response) (result FirewallRuleListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByAccountNextResults retrieves the next set of results, if any.
func (client FirewallRulesClient) listByAccountNextResults(lastResults FirewallRuleListResult) (result FirewallRuleListResult, err error) {
	req, err := lastResults.firewallRuleListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "account.FirewallRulesClient", "listByAccountNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByAccountSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "account.FirewallRulesClient", "listByAccountNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByAccountResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "listByAccountNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByAccountComplete enumerates all values, automatically crossing page boundaries as required.
func (client FirewallRulesClient) ListByAccountComplete(ctx context.Context, resourceGroupName string, accountName string) (result FirewallRuleListResultIterator, err error) {
	result.page, err = client.ListByAccount(ctx, resourceGroupName, accountName)
	return
}

// Update updates the specified firewall rule.
//
// resourceGroupName is the name of the Azure resource group. accountName is the name of the Data Lake Store
// account. firewallRuleName is the name of the firewall rule to update. parameters is parameters supplied to
// update the firewall rule.
func (client FirewallRulesClient) Update(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, parameters *UpdateFirewallRuleParameters) (result FirewallRule, err error) {
	req, err := client.UpdatePreparer(ctx, resourceGroupName, accountName, firewallRuleName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.FirewallRulesClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client FirewallRulesClient) UpdatePreparer(ctx context.Context, resourceGroupName string, accountName string, firewallRuleName string, parameters *UpdateFirewallRuleParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"accountName":       autorest.Encode("path", accountName),
		"firewallRuleName":  autorest.Encode("path", firewallRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2016-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataLakeStore/accounts/{accountName}/firewallRules/{firewallRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	if parameters != nil {
		preparer = autorest.DecoratePreparer(preparer,
			autorest.WithJSON(parameters))
	}
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client FirewallRulesClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client FirewallRulesClient) UpdateResponder(resp *http.Response) (result FirewallRule, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
