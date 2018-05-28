/**
 * Imports
 */
import {Inject, Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import { WiProxyCORSUtils, WiContrib, WiContributionUtils, WiServiceHandlerContribution, AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution, IFieldDefinition, IActionResult, ActionResult, HTTP_METHOD } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";

/**
 * Main
 */
@WiContrib({})
@Injectable()
export class JIRAConnectorContribution extends WiServiceHandlerContribution {
    private category: string;
    constructor(@Inject(Injector) injector, private http: Http) {
        super(injector, http);
        this.category = "Jira";
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;
    }
    
    validate = (name: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
        if (name === "Connect") {
            let name : IFieldDefinition;
            let domain : IFieldDefinition;
            let userName: IFieldDefinition;
            let password: IFieldDefinition;

            for (let configuration of context.settings) {
                if (configuration.name === "name") {
                    name = configuration
                } else if (configuration.name === "domain") {
                    domain = configuration
                } else if (configuration.name === "userName"){
                    userName = configuration
                } else if(configuration.name === "password"){
                    password = configuration
                }
            }
            if (name.value && domain.value && userName.value && password.value) {
                return ValidationResult.newValidationResult().setReadOnly(false)
            } else {
                return ValidationResult.newValidationResult().setReadOnly(true)
            }
        }
        return null;
    }

    action = (actionName: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
        if (actionName == "Connect") {
            return Observable.create(observer => {
                let currentName : string;

                for (let i = 0; i < context.settings.length; i++) {
                    if (context.settings[i].name === "name") {
                        currentName = context.settings[i].value;
                    }
                }

                let duplicate = false;
               
                WiContributionUtils.getConnections(this.http, this.category).subscribe((conns: IConnectorContribution[]) => {
                    for (let conn of conns) {
                        for (let i = 0; i < conn.settings.length; i++) {
                            if (conn.settings[i].name === "name") {
                                let oldName = conn.settings[i].value;
                                if (oldName === currentName && (WiContributionUtils.getUniqueId(conn) !== WiContributionUtils.getUniqueId(context))) {
                                    duplicate = true;
                                    break;
                                }
                            }
                        }
                    }
                
                    if (duplicate) {
                        observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("JIRA-CONNECTOR-001", "Connection name already exists")));
                    } else {
                        let domain = "", userName = "", password = "", project = "";
                        for (let configuration of context.settings) {
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

                        let jiraURL = domain + "/rest/auth/1/session"                        
                        let jiraMetadatURL = domain + "/rest/api/2/issue/createmeta?projectKeys=" + project + "&expand=projects.issuetypes.fields"
                        
                        WiProxyCORSUtils.createRequest(this.http, jiraURL)
                        .addMethod(HTTP_METHOD.GET)
                        .addHeader("Content-Type", "application/json")
                        .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                        .send().subscribe(resp => {
                            console.log("Connection Successful!!1")

                            WiProxyCORSUtils.createRequest(this.http, jiraMetadatURL)
                            .addMethod(HTTP_METHOD.GET)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                            .send().subscribe(resp1 => {
                                console.log("Metadata fetched5")
    
                                if (resp1.json().projects.length > 0) {
                                    let issueTypes = resp1.json().projects[0].issuetypes
                                    var keyvalue: {[key: string]: string} = {}
                                    var notAllowedFields: string[] = ["issuetype","project"];
                                    console.log("abc value ::"+notAllowedFields.indexOf("abc"));
                                    console.log("project index ::"+notAllowedFields.indexOf("project"));
                                    //var keyvalue : Map<string, string> = new Map<string,string>();
                                    //var x: Map<number, string> = new Map<number, string>();
                                    //var keyvalue: Map<string,object> = new Map<string,object>();

                                    for(let i = 0; i <  issueTypes.length; i++) {
                                        if(issueTypes[i].subtask === false) {
                                            let fields = {};                                
                                            let dataType = {};
                                            let jsonRespFields = issueTypes[i].fields;
                                            var requiredFields : string[] = new Array();
                                            for (let fieldName in jsonRespFields) {
                                                if(notAllowedFields.indexOf(fieldName) < 0) {
                                                    dataType["type"] = "string" //assign datatype later discuss
                                                    fields[fieldName] = dataType
                                                    if(jsonRespFields[fieldName].required){
                                                        requiredFields.push(fieldName)
                                                    }
                                                }                                    
                                            }                                
                                        
                                            let inputSchemaForIssueType = {                                    
                                                "properties": fields, "required": requiredFields, "type": "object"
                                            };
                
                                            //console.log("Finish!! inputSchema0 is :: ",JSON.stringify(inputSchemaForIssueType));
                                            keyvalue[issueTypes[i].name] = JSON.stringify(inputSchemaForIssueType)
                                        }
                                    }

                                     //console.log("Value for key defect is ---" +keyvalue["Defect"])
                                     for (let configuration of context.settings) {
                                        if (configuration.name === "keyvalue") {
                                            configuration.value = keyvalue
                                        }
                                    }

                                    let actionResult = {
                                        context: context,
                                        authType: AUTHENTICATION_TYPE.BASIC,
                                        authData: {}
                                    }
                                    observer.next(ActionResult.newActionResult().setSuccess(true).setResult(actionResult));
                                } else {
                                    observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("JIRA-CONNECTOR-003", "Entered Project Key is not valid.")));    
                                }
                            });
                        },
                        error => {
                            observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("JIRA-CONNECTOR-002", "Failed to create connection. Check your configuration.")));
                        });
                    }
                });
            });
        }
    return null;
    }
}