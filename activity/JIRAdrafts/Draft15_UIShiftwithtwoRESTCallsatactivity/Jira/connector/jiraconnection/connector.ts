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
                        let domain = "", userName = "", password = "";
                        for (let configuration of context.settings) {
                            if (configuration.name === "domain") {
                                domain = configuration.value
                            } else if (configuration.name === "userName") {
                                userName = configuration.value
                            } else if (configuration.name === "password") {
                                password = configuration.value
                            }
                        }

                        let jiraURL = domain + "/rest/auth/1/session"
                        let jiraProjectURL = domain + "/rest/api/2/project"
            
                        WiProxyCORSUtils.createRequest(this.http, jiraURL)
                            .addMethod(HTTP_METHOD.GET)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                            .send().subscribe(resp => {
                                
                                console.log("Connection Successful!!");

                                WiProxyCORSUtils.createRequest(this.http, jiraProjectURL)
                                .addMethod(HTTP_METHOD.GET)
                                .addHeader("Content-Type", "application/json")
                                .addHeader("Authorization", "Basic " + btoa(userName + ":" + password))
                                .send().subscribe(resp1 => {
                                    console.log("Projects fetched");
                                    
                                    if (resp1.json().length > 0) {
                                        let projectkeys : string[] = new Array();
                                        for(let i = 0; i < resp1.json().length; i++) {
                                            projectkeys.push(resp1.json()[i].key)                                            
                                        }
                                        
                                        for (let configuration of context.settings) {
                                            if (configuration.name === "projectkeys") {
                                                configuration.value = projectkeys
                                            }
                                        }
    
                                        let actionResult = {
                                            context: context,
                                            authType: AUTHENTICATION_TYPE.BASIC,
                                            authData: {}
                                        }
                                        observer.next(ActionResult.newActionResult().setSuccess(true).setResult(actionResult));
                                    }
                                });
                            },
                            error => {
                                observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("JIRA-CONNECTOR-002", "Failed to create connection. Check your configuration.")));
                            }
                        );
                    }
                });
            });
        }
    return null;
    }
}