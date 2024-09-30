package models

// PlanDetails maps plan names to their detailed descriptions
var PlanDetails = map[string]string{
	"single-region": `Single-Region Plan:

	Description: Access to one probe per minute per target in a single region of your choice.
	Ideal For: Monitoring services critical in a specific geographic location.
	7-Day Free Trial: This plan includes a 7-day free trial period.

Example Usage:
- Monitoring up to 10 targets primarily used by customers in São Paulo, Brazil.
- Selecting the sao-paulo.americas region during sign-up.

Available regions:
Americas:
- Brazil (sao-paulo.americas)
- Canada (canada.americas)
- Canada (calgary.americas)
- United States (northern-virginia.americas)
- United States (ohio.americas)
- United States (northern-california.americas)
- United States (oregon.americas)
Africa:
- South Africa (cape-town.africa)
Asia:
- Japan (tokyo.asia)
- Japan (osaka.asia)
- Hong Kong (hong-kong.asia)
- India (hyderabad.asia)
- India (mumbai.asia)
- Indonesia (jakarta.asia)
- Malaysia (malaysia.asia)
- South Korea (seoul.asia)
- Singapore (singapore.asia)
Oceania:
- Australia (melbourne.oceania)
- Australia (sydney.oceania)
Europe:
- United Kingdom (london.europe)
- Germany (frankfurt.europe)
- Ireland (ireland.europe)
- Italy (milan.europe)
- France (paris.europe)
- Spain (spain.europe)
- Sweden (stockholm.europe)
- Switzerland (zurich.europe)
Middle east:
- Israel (tel-aviv.middle-east)
- Bahrain (bahrain.middle-east)
- United Arab Emirates (uae.middle-east)
`,

	"all-continents": `All-Continents Plan:

Description: Access to one strategically selected region on each continent.
Ideal For: Ensuring global availability and performance across major continents.

Included Regions:
  - Americas:
      - northern-california.americas
      - sao-paulo.americas
  - Africa:
      - cape-town.africa
  - Asia:
      - singapore.asia
  - Oceania:
      - melbourne.oceania
  - Europe:
      - paris.europe
  - Middle East:
      - uae.middle-east

Number of Targets: Select the number of targets you wish to monitor across these regions.`,

	"worldwide": `Worldwide Plan:

Description: Full access to all available regions across the globe.
Ideal For: Comprehensive monitoring for services with a worldwide user base.

Number of Targets: Select the number of targets you wish to monitor across all regions.`,
}
