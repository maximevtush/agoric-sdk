import { makeHelpers } from '@agoric/deploy-script-support';
import { priceFeedProposalBuilder } from './priceFeedSupport.js';

const OPTIONS = {
  AGORIC_INSTANCE_NAME: 'stOSMO-USD price feed',
  IN_BRAND_LOOKUP: ['agoricNames', 'oracleBrand', 'stOSMO'],
  OUT_BRAND_LOOKUP: ['agoricNames', 'oracleBrand', 'USD'],
};

/**
 * @type {import('@agoric/deploy-script-support/src/externalTypes.js').CoreEvalBuilder}
 */
export const defaultProposalBuilder = async ({ publishRef, install }) => {
  return priceFeedProposalBuilder({ publishRef, install }, OPTIONS);
};

export default async (homeP, endowments) => {
  const { writeCoreProposal } = await makeHelpers(homeP, endowments);

  await writeCoreProposal('stOsmoPriceFeed', defaultProposalBuilder);
};
