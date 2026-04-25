package ui

// GfxTheme is the default theme - all other themes derive from this reference
const GfxTheme = `name = "GFX"
description = "GFX default theme - reference for all other themes"

[palette]
# Backgrounds
mainBackgroundColor = "#090D12"       # bunker
inlineBackgroundColor = "#1B2A31"     # dark
selectionBackgroundColor = "#0D141C"  # corbeau

# Text - Content & Body
contentTextColor = "#4E8C93"          # paradiso
labelTextColor = "#8CC9D9"            # dolphin
dimmedTextColor = "#33535B"           # mediterranea
accentTextColor = "#01C2D2"           # caribbeanBlue
highlightTextColor = "#D1D5DA"        # off-white

# Special Text
cwdTextColor = "#67DFEF"              # poseidonJr
footerTextColor = "#519299"           # lagoon

# Borders
boxBorderColor = "#8CC9D9"            # dolphin
separatorColor = "#1B2A31"            # dark

# Confirmation Dialog
confirmationDialogBackground = "#112130"  # trappedDarkness

# UI Elements / Buttons
menuSelectionBackground = "#7EB8C5"   # brighter muted teal
buttonSelectedTextColor = "#0D1418"   # dark text

# Console Output Colors
outputStdoutColor = "#999999"         # neutral gray
outputStderrColor = "#FC704C"         # preciousPersimmon
outputStatusColor = "#01C2D2"         # caribbeanBlue
outputWarningColor = "#F2AB53"        # safflower
outputDebugColor = "#33535B"          # mediterranea
outputInfoColor = "#01C2D2"           # caribbeanBlue
spinnerColor = "#FC704C"              # preciousPersimmon
`

// SpringTheme is a spring-themed color palette with greens and vibrant energy
const SpringTheme = `name = "Spring"
description = "Fresh spring greens with vibrant energy"

[palette]
# Backgrounds - sapphire → ceruleanBlue → sapphire gradient
mainBackgroundColor = "#323B9E"       # sapphire (main background)
inlineBackgroundColor = "#0972BB"     # easternBlue (secondary areas)
selectionBackgroundColor = "#090D12"  # bunker (highlight areas)

# Text - Content & Body - green colors for positive, red for negative
contentTextColor = "#179CA8"          # easternBlue - neutral readable
labelTextColor = "#90D88D"            # feijoa (labels)
dimmedTextColor = "#C8E189"           # yellowGreen (dimmed)
accentTextColor = "#FEEA85"           # salomie - bright shortcuts
highlightTextColor = "#D1D5DA"        # off-white (highlights)

# Special Text
cwdTextColor = "#FEEA85"              # salomie - bright yellow accent
footerTextColor = "#58C9BA"           # downy - muted descriptions

# Borders
boxBorderColor = "#90D88D"            # feijoa
separatorColor = "#0972BB"            # easternBlue

# Confirmation Dialog
confirmationDialogBackground = "#244DA8"  # ceruleanBlue

# UI Elements
menuSelectionBackground = "#5BCF90"   # emerald - natural green
buttonSelectedTextColor = "#3F2894"   # daisyBush - dark contrast

# Console Output Colors
outputStdoutColor = "#999999"         # neutral gray
outputStderrColor = "#FD5B68"         # wildWatermelon
outputStatusColor = "#4ECB71"         # emerald
outputWarningColor = "#F67F78"        # froly
outputDebugColor = "#C8E189"          # yellowGreen
outputInfoColor = "#37CB9F"           # shamrock
spinnerColor = "#FD5B68"              # wildWatermelon
`

// SummerTheme is a summer-themed color palette with electric blues and bright sunshine
const SummerTheme = `name = "Summer"
description = "Warm summer blues and bright sunshine"

[palette]
# Backgrounds - blueMarguerite → havelockBlue → violetBlue
mainBackgroundColor = "#000000"       # black (main background)
inlineBackgroundColor = "#4D88D1"     # havelockBlue (secondary areas)
selectionBackgroundColor = "#090D12"  # bunker (highlight areas)

# Text - Content & Body - electric cyan/bright for positives, hot reds for negatives
contentTextColor = "#3CA7E0"          # violetBlue - readable neutral
labelTextColor = "#19E5FF"            # cyan (labels)
dimmedTextColor = "#5E68C1"           # indigo (dimmed)
accentTextColor = "#FFBF16"           # lightningYellow - electric shortcuts
highlightTextColor = "#D1D5DA"        # off-white (highlights)

# Special Text
cwdTextColor = "#FFBF16"              # lightningYellow - electric accent
footerTextColor = "#8667BF"           # blueMarguerite - muted descriptions

# Borders
boxBorderColor = "#19E5FF"            # cyan
separatorColor = "#4D88D1"            # havelockBlue

# Confirmation Dialog
confirmationDialogBackground = "#2BC6F0"  # pictonBlue

# UI Elements
menuSelectionBackground = "#FE62B9"   # hotPink - electric highlight
buttonSelectedTextColor = "#8667BF"   # blueMarguerite - dark contrast

# Console Output Colors
outputStdoutColor = "#999999"         # neutral gray
outputStderrColor = "#FF3469"         # radicalRed
outputStatusColor = "#00FFFF"         # electric cyan
outputWarningColor = "#FF9700"        # pizazz
outputDebugColor = "#5E68C1"          # indigo
outputInfoColor = "#2BC6F0"           # pictonBlue
spinnerColor = "#FF3469"              # radicalRed
`

// AutumnTheme is an autumn-themed color palette with rich golds and warm earth tones
const AutumnTheme = `name = "Autumn"
description = "Rich autumn oranges and warm earth tones"

[palette]
# Backgrounds - jacaranda → mulberryWood → roseBudCherry
mainBackgroundColor = "#3E0338"       # jacaranda (main background)
inlineBackgroundColor = "#5E063E"     # mulberryWood (secondary areas)
selectionBackgroundColor = "#090D12"  # bunker (highlight areas)

# Text - Content & Body - gold colors for positive, deep reds for negative
contentTextColor = "#E78C79"          # apricot - warm readable
labelTextColor = "#F9C94D"            # saffronMango (labels)
dimmedTextColor = "#F09D06"           # tulipTree (dimmed)
accentTextColor = "#F5BB09"           # corn - bright shortcuts
highlightTextColor = "#D1D5DA"        # off-white (highlights)

# Special Text
cwdTextColor = "#F5BB09"              # corn - golden bright
footerTextColor = "#CD5861"           # chestnutRose - muted descriptions

# Borders
boxBorderColor = "#F9C94D"            # saffronMango
separatorColor = "#5E063E"            # mulberryWood

# Confirmation Dialog
confirmationDialogBackground = "#7D0E36"  # roseBudCherry

# UI Elements
menuSelectionBackground = "#F1AE37"   # tulipTree - golden harvest highlight
buttonSelectedTextColor = "#3E0338"   # jacaranda - darkest contrast

# Console Output Colors
outputStdoutColor = "#999999"         # neutral gray
outputStderrColor = "#DC3003"         # grenadier
outputStatusColor = "#F5BB09"         # corn
outputWarningColor = "#E85C03"        # trinidad
outputDebugColor = "#F09D06"          # tulipTree
outputInfoColor = "#F48C06"           # tangerine
spinnerColor = "#F9C94D"              # saffronMango
`

// WinterTheme is a winter-themed color palette with professional blues and subtle elegance
const WinterTheme = `name = "Winter"
description = "Cool winter purples with subtle elegance"

[palette]
# Backgrounds - cloudBurst → sanJuan → sanMarino
mainBackgroundColor = "#233253"       # cloudBurst (main background)
inlineBackgroundColor = "#334676"     # sanJuan (secondary areas)
selectionBackgroundColor = "#090D12"  # bunker (highlight areas)

# Text - Content & Body - professional blues for positive, soft pinks for negative
contentTextColor = "#CAD0E6"          # cyanGray - cool readable
labelTextColor = "#7F95D6"            # chetwodeBlue (labels)
dimmedTextColor = "#9BA9D0"           # rockBlue (dimmed)
accentTextColor = "#F6F5FA"           # whisper - bright shortcuts
highlightTextColor = "#D1D5DA"        # off-white (highlights)

# Special Text
cwdTextColor = "#F6F5FA"              # whisper - bright white
footerTextColor = "#9BA9D0"           # rockBlue - muted descriptions

# Borders
boxBorderColor = "#7F95D6"            # chetwodeBlue
separatorColor = "#334676"            # sanJuan

# Confirmation Dialog
confirmationDialogBackground = "#233253"  # cloudBurst

# UI Elements
menuSelectionBackground = "#7F95D6"   # chetwodeBlue - professional blue accent
buttonSelectedTextColor = "#F6F5FA"   # whisper - light contrast

# Console Output Colors
outputStdoutColor = "#999999"         # neutral gray
outputStderrColor = "#E0BACF"         # melanie
outputStatusColor = "#435A98"         # sanMarino
outputWarningColor = "#CEBAC5"        # lily
outputDebugColor = "#9BA9D0"          # rockBlue
outputInfoColor = "#435A98"           # sanMarino
spinnerColor = "#F6F5FA"              # whisper
`
