/**
 * Imports
 */
import {Inject, Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import { WiProxyCORSUtils, WiContrib, WiContributionUtils, WiServiceHandlerContribution, AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution, IFieldDefinition, IActionResult, ActionResult, HTTP_METHOD } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { Util, Connection } from "./connector.util"

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
            let conn: Connection = Util.getConnection(context.settings)
            if (conn.name && conn.domain && conn.userName && conn.password) {
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
                WiContributionUtils.getConnections(this.http, this.category).subscribe((configurations: IConnectorContribution[]) => {
                    var conn: Connection = Util.getConnection(context.settings)
                    if (Util.isDuplicate(context, configurations)) {
                        observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("JIRA-CONNECTOR-1001", "Connection name already exists")));
                    } else {
                        let jiraURL = conn.domain + "/rest/auth/1/session"
            
                        WiProxyCORSUtils.createRequest(this.http, jiraURL)
                            .addMethod(HTTP_METHOD.GET)
                            .addHeader("Content-Type", "application/json")
                            .addHeader("Authorization", "Basic " + btoa(conn.userName + ":" + conn.password))
                            .send().subscribe(resp => {
                                console.log("Connection Successful!!");
                                let actionResult = {
                                    context: context,
                                    authType: AUTHENTICATION_TYPE.BASIC,
                                    authData: {}
                                }
                                observer.next(ActionResult.newActionResult().setSuccess(true).setResult(actionResult));
                        },
                        error => {
                            observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("JIRA-CONNECTOR-1002", "Failed to create connection. Check your configuration.")));
                        });
                    }
                });
            });
        }
    return null;
    }
}