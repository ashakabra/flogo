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

export class JiraQueryTicketActivityContribution extends WiServiceHandlerContribution {
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
        } else if (fieldName === "issueTypefromConn") {
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
    
                            console.log("Domain is :"+domain+" userName is :"+userName);
                            
                            //VERY IMPORTANT
                            //https://jira.tibco.com/rest/api/2/project/FILB this rest api need to be used
                            //let jiraIssueTypeURL = domain + "/rest/api/2/issuetype"
                            let projectName = "FILB"
                            let jiraIssueTypeURL = domain + "/rest/api/2/project/" + projectName
                            WiProxyCORSUtils.createRequest(this.http, jiraIssueTypeURL)
                                .addMethod(HTTP_METHOD.GET)
                                .addHeader("Content-Type", "application/json")
                                .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                                .send().subscribe(resp => {    
                                    console.log("Issue types Fetched1 :: "+resp.json().issueTypes.length);                                
                                    //console.log("ASHa5 json response = " + JSON.stringify(resp.json()));
                                    var issueTypefromConnRefs : string[] = new Array();
                                    /*for(let i = 0; i < resp.json().length; i++){
                                        console.log("ASHA8 json response = " + resp.json()[0].subtask);
                                        if(resp.json()[i].subtask === "true" || resp.json()[i].description === ""){
                                            console.log("ASHA8 json response = " + resp.json()[i].name);
                                        } else {
                                            issueTypefromConnRefs.push(resp.json()[i].name);
                                        }
                                    }*/
                                    for(let i = 0; i < resp.json().issueTypes.length; i++){
                                        console.log("ASHA8 json response = " + resp.json().issueTypes[i].name);
                                        /*if(resp.json()[i].subtask === "true" || resp.json()[i].description === ""){
                                            console.log("ASHA8 json response = " + resp.json()[i].name);
                                        } else {*/
                                            //issueTypefromConnRefs.push(resp.json()[i].name);
                                        //}
                                    }
                                    observer.next(issueTypefromConnRefs);
                                },
                                error => {
                                    console.log("Failed to get fields");
                                    observer.next("{}");
                                }
                            );                              
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