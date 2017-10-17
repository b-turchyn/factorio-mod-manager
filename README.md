# Factorio Mod Manager

A modpack manager for Factorio

## Installation and Usage

Create a folder named `modpacks` in your Factorio UserData directory (typically
%APPDATA%\Factorio). Create a subfolder for each modpack you have, and put the
mods for that modpack in each subfolder. 

## Usage

Double-click on `factorio-mod-manager.exe`, and select the number of the modpack
you wish to apply.

From command line (non-interactive): 

    factorio-mod-manager.exe <modpack>

To remove all mods, simply run with a modpack of `Base`, or when running
interactively, use modpack 0, as shown in the UI. 
    
# Caveats

* This mod manager wipes out the contents of the `mods` directory when running.
  Assume that everything you're running will be removed prior to a modpack
  being applied.
* The currently-selected modpack is not tracked.
* Updates applied in the `mods` directory need to be manually copied back to the
  associated modpack directory.
