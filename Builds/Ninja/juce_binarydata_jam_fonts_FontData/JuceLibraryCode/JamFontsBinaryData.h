/* =========================================================================================

   This is an auto-generated file: Any edits you make may be overwritten!

*/

#pragma once

namespace jam::fonts
{
    extern const char*   DisplayBold_ttf;
    const int            DisplayBold_ttfSize = 56212;

    extern const char*   DisplayBook_ttf;
    const int            DisplayBook_ttfSize = 58168;

    extern const char*   DisplayMedium_ttf;
    const int            DisplayMedium_ttfSize = 58492;

    extern const char*   DisplayMonoBold_ttf;
    const int            DisplayMonoBold_ttfSize = 237584;

    extern const char*   DisplayMonoBook_ttf;
    const int            DisplayMonoBook_ttfSize = 240056;

    extern const char*   DisplayMonoMedium_ttf;
    const int            DisplayMonoMedium_ttfSize = 239548;

    // Number of elements in the namedResourceList and originalFileNames arrays.
    const int namedResourceListSize = 6;

    // Points to the start of a list of resource names.
    extern const char* namedResourceList[];

    // Points to the start of a list of resource filenames.
    extern const char* originalFilenames[];

    // If you provide the name of one of the binary resource variables above, this function will
    // return the corresponding data and its size (or a null pointer if the name isn't found).
    const char* getNamedResource (const char* resourceNameUTF8, int& dataSizeInBytes);

    // If you provide the name of one of the binary resource variables above, this function will
    // return the corresponding original, non-mangled filename (or a null pointer if the name isn't found).
    const char* getNamedResourceOriginalFilename (const char* resourceNameUTF8);
}
