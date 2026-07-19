import type React from "react";

export interface CommercialComponents {
	OidcConfig: React.ComponentType<any> | null;
	RbacManager: React.ComponentType<any> | null;
	AuditLogViewer: React.ComponentType<any> | null;
}

export const commercialRegistry: CommercialComponents = {
	OidcConfig: null,
	RbacManager: null,
	AuditLogViewer: null,
};
