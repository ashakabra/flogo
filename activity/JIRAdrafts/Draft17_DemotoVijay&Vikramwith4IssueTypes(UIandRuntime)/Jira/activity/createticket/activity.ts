/**
 * Imports
 */
import { Observable } from "rxjs/Observable";
import { Injectable, Injector, Inject } from "@angular/core";
import { Http } from "@angular/http";
import {
    WiContrib,
    WiProxyCORSUtils,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IFieldDefinition,
    IActivityContribution,
    IConnectorContribution,
    WiContributionUtils
} from "wi-studio/app/contrib/wi-contrib";
import {HTTP_METHOD} from "wi-studio/common/models/contrib";
import { JsonSchema } from "./activity.jsonschema";

export class Project {
    name: string;
    issueType: IssueTypes[]; 
}
export class IssueTypes {
    name: string;    
}

@WiContrib({})
@Injectable()

export class JiraCreateTicketActivityContribution extends WiServiceHandlerContribution {
    private category: string;
    constructor( @Inject(Injector) injector, private http: Http) {
        super(injector, http);
        this.category = "Jira";
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "Connection") {
            return Observable.create(observer => {
                let connectionRefs = [];
                
                WiContributionUtils.getConnections(this.http, this.category).subscribe((data: IConnectorContribution[]) => {
                    data.forEach(connection => {
                        for (let i = 0; i < connection.settings.length; i++) {
                            if (connection.settings[i].name === "name") {
                                connectionRefs.push({
                                    "unique_id": WiContributionUtils.getUniqueId(connection),
                                    "name": connection.settings[i].value
                                });
                                break;
                            }
                        }
                    });
                    observer.next(connectionRefs);
                });
            });
        } else if (fieldName === "project") {
            let connectionField: IFieldDefinition = context.getField("Connection");
           
            if(connectionField.value) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectionField.value)
                    .map(data => data)
                    .subscribe(data =>{
                        var projectKeys : string[] = new Array();                      
                        for (let configuration of data.settings) {
                            if (configuration.name === "projandissue") {
                                var keyval = JSON.stringify(configuration.value)
                                let prjs:Project[] = JSON.parse(keyval)
                                for (let prj of prjs) {
                                    projectKeys.push(prj.name)
                                }
                                observer.next(projectKeys);
                            }
                        }                                                                      
                    });
                });
            }
        } else if (fieldName === "issueType") {
            let connectionField: IFieldDefinition = context.getField("Connection");
            let projectField: IFieldDefinition = context.getField("project");

            if(connectionField.value && projectField.value) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectionField.value)
                    .map(data => data)
                    .subscribe(data =>{
                        var issueTypes : string[] = new Array();
                        for (let configuration of data.settings) {
                            if (configuration.name === "projandissue") {
                                var keyval = JSON.stringify(configuration.value)
                                let prjs:Project[] = JSON.parse(keyval)
                                for (let prj of prjs) {
                                    if(prj.name == projectField.value) {
                                        for (let issues of prj.issueType) {
                                            issueTypes.push(issues.name)
                                        }
                                        break;
                                    }
                                }
                                observer.next(issueTypes);
                            } 
                        }
                      
                                                    
                    });
                });
            
            }
        } else if (fieldName === "input" || fieldName === "output") {
            let connectionField: IFieldDefinition = context.getField("Connection");
            let projectField: IFieldDefinition = context.getField("project");
            let issueTypeField: IFieldDefinition = context.getField("issueType");

            if(connectionField.value && projectField.value && issueTypeField.value){
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectionField.value)
                    .map(data => data)
                    .subscribe(data =>{
                        let domain = "", userName = "", password = "";
                        for (let configuration of data.settings) {
                            if (configuration.name === "domain") {
                                domain = configuration.value
                            } else if (configuration.name === "userName") {
                                userName = configuration.value
                            } else if (configuration.name === "password") {
                                password = configuration.value
                            }
                        }
        
                        let jiraMetadataURL = domain + "/rest/api/2/issue/createmeta?projectKeys="+projectField.value+"&issuetypeNames="+issueTypeField.value+"&expand=projects.issuetypes.fields";
                        console.log("Before Metadata REST API!!")
                        WiProxyCORSUtils.createRequest(this.http, jiraMetadataURL)
                            .addMethod(HTTP_METHOD.GET)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                            .send().subscribe(resp => {
                                console.log("Successful result from Metadata REST API!!");
                                                                
                                let notAllowedFields: string[] = ["issuetype","project","issuelinks","timetracking","attachment"];
                                let reqFields = {}, optFields = {};
                                let dataType = {};
                                let jsonRespFields = resp.json().projects[0].issuetypes[0].fields;
                                var requiredFields : string[] = new Array();
                                for (let jiraFieldName in jsonRespFields) {
                                    if(notAllowedFields.indexOf(jiraFieldName) < 0) {                                       
                                        if(jsonRespFields[jiraFieldName].required){
                                            reqFields[jsonRespFields[jiraFieldName].name] = JsonSchema.Types.toJsonType(jsonRespFields[jiraFieldName].schema, jiraFieldName )
                                            requiredFields.push(jsonRespFields[jiraFieldName].name)
                                        } else {
                                            optFields[jsonRespFields[jiraFieldName].name] = JsonSchema.Types.toJsonType(jsonRespFields[jiraFieldName].schema, jiraFieldName)
                                        }
                                        var fields = { ...reqFields, ...optFields };
                                        if (fieldName === "output") {
                                            let outputKey = { "IssueID" : {
                                                "type":"string"
                                            }}
                                            fields = { ...outputKey, ...fields}
                                        }                                                 
                                    }                              
                                }                                
                                let inputSchema = {                                    
                                    "properties": fields, "required": requiredFields, "type": "object"
                                };                                
                                //console.log("Finish!! inputSchema0 is :: ",JSON.stringify(inputSchema));
                                observer.next(JSON.stringify(inputSchema));
                            },
                            error => {
                                console.log("Failed to get fields");
                                observer.next("{}");
                            });
                    });
                });
            }
        }
        return null;
    }
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "Connection") {
            let connection: IFieldDefinition = context.getField("Connection")
            if (connection.value === null) {
                return ValidationResult.newValidationResult().setError("JIRA-1000", "Jira Connection must be configured");
            }
        }
        return null;
    }
}