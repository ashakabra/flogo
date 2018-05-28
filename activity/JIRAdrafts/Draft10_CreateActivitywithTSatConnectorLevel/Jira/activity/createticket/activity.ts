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
        } else if (fieldName === "issueType") {
            let connectionField: IFieldDefinition = context.getField("Connection");            

            if(connectionField.value) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectionField.value)
                    .map(data => data)
                    .subscribe(data =>{
                        let domain = "", userName = "", password = "", project = "";
                        for (let configuration of data.settings) {
                            if (configuration.name === "domain") {
                                domain = configuration.value
                            } else if (configuration.name === "userName") {
                                userName = configuration.value
                            } else if (configuration.name === "password") {
                                password = configuration.value
                            } else if (configuration.name === "project") {
                                project = configuration.value
                            }
                        }
                      
                        let jiraIssueTypeURL = domain + "/rest/api/2/project/" + project
                        console.log("Before IssueType REST API!!");
                        WiProxyCORSUtils.createRequest(this.http, jiraIssueTypeURL)
                            .addMethod(HTTP_METHOD.GET)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                            .send().subscribe(resp => {  
                                console.log("Successful result from IssueType REST API!!");            
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
            let issueTypeField: IFieldDefinition = context.getField("issueType");
            
            if(connectionField.value && issueTypeField.value) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectionField.value)
                    .map(data => data)
                    .subscribe(data =>{
                        
                        var keyvaluefromConn: {[key: string]: string} = {}
                        for (let configuration of data.settings) {
                            if (configuration.name === "keyvalue") {
                                keyvaluefromConn = configuration.value
                            } 
                        }
                        observer.next(keyvaluefromConn[issueTypeField.value]);
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