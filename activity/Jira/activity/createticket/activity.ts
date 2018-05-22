/**
 * Imports
 */
import { Observable } from "rxjs/Observable";
import { Injectable, Injector, Inject } from "@angular/core";
import { Http } from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IFieldDefinition,
    IActivityContribution,
    IConnectorContribution,
    WiContributionUtils
} from "wi-studio/app/contrib/wi-contrib";

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
        }
        return null;
    }
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "Connection") {
            let connection: IFieldDefinition = context.getField("Connection")
            if (connection.value === null) {
                return ValidationResult.newValidationResult().setError("JIRA-1000", "Jira Connection must be configured");
            }
        } else if (fieldName === "severity") {
            let issue: IFieldDefinition = context.getField("issueType");
            if (issue.value === "Task") {
                return ValidationResult.newValidationResult().setVisible(false);
            } else {
                return ValidationResult.newValidationResult().setVisible(true);
            }
        } else if (fieldName === "confirmer") {
            let issue: IFieldDefinition = context.getField("issueType");
            if (issue.value === "Story" || issue.value === "Task") {
                return ValidationResult.newValidationResult().setVisible(false);
            } else {
                return ValidationResult.newValidationResult().setVisible(true);
            }
        } else if (fieldName === "affectVersion") {
            let issue: IFieldDefinition = context.getField("issueType");
            if (issue.value === "Task") {
                return ValidationResult.newValidationResult().setVisible(false);
            } else {
                return ValidationResult.newValidationResult().setVisible(true);
            }
        }
        return null;
    }
}