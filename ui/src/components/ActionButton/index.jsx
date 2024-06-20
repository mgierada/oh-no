"use client";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import formatCurrentDate from "@/utils/helper";

/**
 * @typedef {Object} ActionButtonProps
 * @property {string} toast_message - the message to be displayed on on the toast popup
 * @property {string} cta_button - the call to action to be displayed on the button
 */

/**
 * ActionButton component
 * @param {ActionButtonProps} actionButtonProps
 * @returns {JSX.Element}
 */
export function ActionButton({ toast_message, cta_button }) {
  return (
    <Button
      variant="outline"
      onClick={() =>
        toast(
          toast_message,
          {
            description: formatCurrentDate(),
            action: {
              label: "Ok",
              // onClick: () => console.log("Undo"),
            },
          },
          console.log("Event has been recorded"),
        )
      }
    >
      {cta_button}
    </Button>
  );
}
