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
                                                        
                        let jiraProjectURL = domain + "/rest/api/2/project"   
                        console.log("Before Project REST API!!");                         
                        WiProxyCORSUtils.createRequest(this.http, jiraProjectURL)
                        .addMethod(HTTP_METHOD.GET)
                        .addHeader("Content-Type", "application/json")
                        .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                        .send().subscribe(resp => {    
                            console.log("Successful result from Project REST API!!");                         
                            var projectkeys : string[] = new Array();
                                for(let i = 0; i < resp.json().length; i++){
                                    projectkeys.push(resp.json()[i].key);                                        
                                }
                                observer.next(projectkeys);
                            },
                            error => {
                                console.log("Failed to get fields");
                                observer.next("{}");
                            }
                        );                              
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
                      
                        let jiraIssueTypeURL = domain + "/rest/api/2/project/" + projectField.value
                        console.log("Before IssueType REST API!!");
                        WiProxyCORSUtils.createRequest(this.http, jiraIssueTypeURL)
                            .addMethod(HTTP_METHOD.GET)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                            .send().subscribe(resp => {  
                                console.log("Successful result from IssueType REST API!!");  
                                //console.log("ASHa5 json response = " + JSON.stringify(resp.json()));
                                var issueTypes : string[] = new Array();
                                for(let i = 0; i < resp.json().issueTypes.length; i++) {
                                    if(resp.json().issueTypes[i].subtask === false) {
                                        issueTypes.push(resp.json().issueTypes[i].name)
                                    }
                                }
                                observer.next(issueTypes);
                            },
                            error => {
                                console.log("Failed to get fields");
                                observer.next("{}");
                            }
                        );                              
                    });
                });
            
            }
        } else if (fieldName === "input") {
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
                                
                                //let type = {"type":"string"}
                                let fields = {};                                
                                let dataType = {};
                                let jsonRespFields = resp.json().projects[0].issuetypes[0].fields;
                                var requiredFields : string[] = new Array();
                                for (let fieldName in jsonRespFields) {                                    
                                    dataType["type"] = "string" //assign datatype later discuss
                                    //fields[fieldName] = type
                                    fields[fieldName] = dataType
                                    if(jsonRespFields[fieldName].required){
                                        requiredFields.push(fieldName)
                                    }                                    
                                }                                
                                let inputSchema = {                                    
                                    "properties": fields, "required": requiredFields, "type": "object"
                                };
                                console.log("Finish!! inputSchema0 is :: ",JSON.stringify(inputSchema));
                                observer.next(JSON.stringify(inputSchema));
                                //console.log("Data1 is :: "+resp.json().projects[0].issuetypes[0].fields.issuetype.required);                                
                            },
                            error => {
                                console.log("Failed to get fields");
                                observer.next("{}");
                            });
                            
                            //let json = {"abc":"string","pqr":"string"}
                            //observer.next(JSON.stringify(json));
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