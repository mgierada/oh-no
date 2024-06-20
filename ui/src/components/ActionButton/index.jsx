"use client";
import { toast } from "sonner";
import { Biohazard, Activity } from "lucide-react";

import { Button } from "@/components/ui/button";
import getAndFormatCurrentDate from "@/utils/helper";

/**
 * @typedef {Object} ActionButtonProps
 * @property {string} variant - a string that determines the variant of the button
 * @property {string} toast_message - the message to be displayed on on the toast popup
 * @property {string} cta_button - the call to action to be displayed on the button
 * @property {string} icon - the icon to be displayed on the button
 */

/**
 * ActionButton component
 * @param {ActionButtonProps} actionButtonProps
 * @returns {JSX.Element}
 */
export function ActionButton({ variant, toast_message, cta_button, icon }) {
  const renderIcon = (icon) => {
    if (icon === "biohazard") {
      return <Biohazard className="mr-2 h-4 w-4" />;
    } else if (icon === "activity") {
      return <Activity className="mr-2 h-4 w-4" />;
    }
    return null;
  };
  return (
    <Button
      variant={variant}
      onClick={() =>
        toast(
          toast_message,
          {
            description: getAndFormatCurrentDate(),
            action: {
              label: "Ok",
              // onClick: () => console.log("Undo"),
            },
          },
          console.log("Event has been recorded"),
        )
      }
    >
      {renderIcon(icon)} {cta_button}
    </Button>
  );
}
