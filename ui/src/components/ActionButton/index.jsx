"use client";
import { useState } from "react";
import { toast } from "sonner";
import { Biohazard, Activity } from "lucide-react";

import { Button } from "@/components/ui/button";
import getAndFormatCurrentDate from "@/utils/helper";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from "@/components/ui/alert-dialog";

/**
 * @typedef {Object} ActionButtonProps
 * @property {string} variant - a string that determines the variant of the button
 * @property {string} toastMessage - the message to be displayed on the toast popup
 * @property {string} ctaButton - the call to action to be displayed on the button
 * @property {string} icon - the icon to be displayed on the button
 * @property {string} alertDialogDescription - description text shown on the alert dialog
 * @property {function} handleUpdate - function to be called when the button is clicked
 */

/**
 * ActionButton component
 * @param {ActionButtonProps} actionButtonProps
 * @returns {JSX.Element}
 */
export function ActionButton({
  variant,
  toastMessage,
  ctaButton,
  icon,
  alertDialogDescription,
  handleUpdate,
}) {
  const [showAlertDialog, setShowAlertDialog] = useState(false);

  const renderIcon = (icon) => {
    if (icon === "biohazard") {
      return <Biohazard className="mr-2 h-4 w-4" />;
    } else if (icon === "activity") {
      return <Activity className="mr-2 h-4 w-4" />;
    }
    return null;
  };

  const handleConfirm = async () => {
    try {
      await handleUpdate();
      toast(toastMessage, {
        description: getAndFormatCurrentDate(),
        action: {
          label: "Dismiss",
        },
      });
      console.log("Event has been recorded");
      setShowAlertDialog(false);
    } catch (error) {
      console.error("Failed to record an event:", error);
      toast("Failed to record an event", {
        description: error.message,
        action: {
          label: "Dismiss",
        },
      });
    }
  };

  return (
    <>
      <Button variant={variant} onClick={() => setShowAlertDialog(true)}>
        {renderIcon(icon)} {ctaButton}
      </Button>

      {showAlertDialog && (
        <AlertDialog open={showAlertDialog} onOpenChange={setShowAlertDialog}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure? 🤔</AlertDialogTitle>
              <AlertDialogDescription>
                {alertDialogDescription}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel onClick={() => setShowAlertDialog(false)}>
                Cancel
              </AlertDialogCancel>
              <AlertDialogAction onClick={handleConfirm}>
                Sure, let's go!
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </>
  );
}
