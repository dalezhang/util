class Biz::CsgoMatchInfoUpdater
  def initialize(data, logger = nil)
  	logger ||= ::Biz::MarkLogger.new
    @logger = logger
    @logger.mark("Biz::CsgoMatchInfoUpdater#initialize")
    @data = data
    @extern_id = @data["listId"]
  end

  def run
    @logger ||= ::Biz::MarkLogger.new
    @logger.mark("Biz::CsgoMatchInfoUpdater#run")
    if !@extern_id.present?
      @logger.log_info("!@extern_id.present?")
      return
    end
    db_series = CsgoSeries.find_by(extern_id: @extern_id)
    db_series = CsgoSeries.find_by(extern_id: "hltv_#{@extern_id}") if db_series.nil?
    if !db_series.present?
      @logger.log_info("!db_series.present? @extern_id = #{@extern_id}")
      return
    end
    series_left_score = 0
    series_right_score = 0
    @data["mapScores"].each do |key, value|
      match = db_series.matches.find_by(game_no: key)
      if !match.present?
        @logger.log_info("!match.present? game_no = #{key}")
        next
      end
      value["scores"].each do |k,v|
        team = Team.find_by(extern_id: k)
        next if !team.present?
        if team.id == db_series.left_team_id
          match.update!(left_score: v)
        end
        if team.id == db_series.right_team_id
          match.update!(right_score: v)
        end
      end
      if value["mapOver"] 
        if match.left_score > match.right_score
          series_left_score += 1
        end
        if match.left_score < match.right_score
          series_right_score += 1
        end
      end
      score_info_hash = {}
      if value["firstHalf"].present?
        score_info_hash[:first_half] = {}
        score_info_hash[:first_half][:ct_score] = value["firstHalf"]["ctScore"]
        score_info_hash[:first_half][:t_score] = value["firstHalf"]["tScore"]
        ct_team = Team.find_by(extern_id: value["firstHalf"]["ctTeamDbId"])
        # t_team = Team.find_by(extern_id: value["firstHalf"]["tTeamDbId"])
        if ct_team.present? #&& t_team.present?
          if ct_team.id == db_series.left_team_id
            score_info_hash[:first_half][:ct] = "left"
            score_info_hash[:first_half][:t] = "right"
          end
          if ct_team.id == db_series.right_team_id
            score_info_hash[:first_half][:ct] = "right"
            score_info_hash[:first_half][:t] = "left"
          end
        end
      end
      if value["secondHalf"].present?
        score_info_hash[:second_half] = {}
        score_info_hash[:second_half][:ct_score] = value["secondHalf"]["ctScore"]
        score_info_hash[:second_half][:t_score] = value["secondHalf"]["tScore"]
        ct_team = Team.find_by(extern_id: value["secondHalf"]["ctTeamDbId"])
        # t_team = Team.find_by(extern_id: value["secondHalf"]["tTeamDbId"])
        if ct_team.present? #&& t_team.present?
          if ct_team.id == db_series.left_team_id
            score_info_hash[:second_half][:ct] = "left"
            score_info_hash[:second_half][:t] = "right"
          end
          if ct_team.id == db_series.right_team_id
            score_info_hash[:second_half][:ct] = "right"
            score_info_hash[:second_half][:t] = "left"
          end
        end
      end
      if value["overtime"].present?
        score_info_hash[:over_time] = {}
        score_info_hash[:over_time][:ct_score] = value["overtime"]["ctScore"]
        score_info_hash[:over_time][:t_score] = value["overtime"]["tScore"]
        ct_team = Team.find_by(extern_id: value["overtime"]["ctTeamDbId"])
        # t_team = Team.find_by(extern_id: value["secondHalf"]["tTeamDbId"])
        if ct_team.present? #&& t_team.present?
          if ct_team.id == db_series.left_team_id
            score_info_hash[:over_time][:ct] = "left"
            score_info_hash[:over_time][:t] = "right"
          end
          if ct_team.id == db_series.right_team_id
            score_info_hash[:over_time][:ct] = "right"
            score_info_hash[:over_time][:t] = "left"
          end
        end
      end
      if score_info_hash.present?
        match.update!(score_info: score_info_hash.to_json)
        @logger.mark("match #{key} update score_info #{score_info_hash}")
      end
      if value["map"].present?
        map_name = value["map"].split("_")[1].capitalize rescue nil
        match.update!(map: map_name) if map_name.present?
        @logger.mark("match #{key} update map #{map_name}")
      end
      @logger.mark("match #{key} update scores #{value["scores"]}")
    end
    @logger.log_info("series #{@extern_id} data: \n #{@data}")
    db_series.update!(left_score: series_left_score, right_score: series_right_score)
  rescue => e
    raise @logger.error(e)
  end

  def self.subscribe(extern_id, logger = nil)
  	logger ||= ::Biz::MarkLogger.new
    @logger ||= logger
    @logger.mark("Biz::CsgoMatchInfoUpdater#subscribe")
    url = ENV["CSGO_MATCH_INFO_SUBSCRIBE_URL"]
    if !url.present?
      raise "CSGO_MATCH_INFO_SUBSCRIBE_URL is nil"
    end
    response = HTTParty.get("#{url}?match_id=#{extern_id}")
    if response.code > 300
      @logger.log_info("extern_id = #{extern_id}, response.body = \n#{response.body}")
    end
    @logger.log_info("subscribe #{extern_id}")
  rescue => e
    raise @logger.error(e)
  end

  def self.upsubscribe(extern_id, logger = nil)
  	logger ||= ::Biz::MarkLogger.new
    @logger ||= logger
    @logger.mark("Biz::CsgoMatchInfoUpdater#subscribe")
    url = ENV["CSGO_MATCH_INFO_SUBSCRIBE_URL"]
    if !url.present?
      raise "CSGO_MATCH_INFO_SUBSCRIBE_URL is nil"
    end
    response = HTTParty.delete("#{url}?match_id=#{extern_id}")
    if response.code > 300
      @logger.log_info("extern_id = #{extern_id}, response.body = \n#{response.body}")
    end
    @logger.log_info("upsubscribe #{extern_id}")
  rescue => e
    raise @logger.error(e)
  end
end

#  {"mapScores"=>
#   {"1"=>
#     {"firstHalf"=>{"ctTeamDbId"=>6010, "ctScore"=>3, "tTeamDbId"=>5293, "tScore"=>12},
#      "secondHalf"=>{"ctTeamDbId"=>5293, "ctScore"=>3, "tTeamDbId"=>6010, "tScore"=>0},
#      "overtime"=>{"ctTeamDbId"=>5293, "ctScore"=>0, "tTeamDbId"=>6010, "tScore"=>0},
#      "mapOrdinal"=>1,
#      "scores"=>{"5293"=>15, "6010"=>3},
#      "currentCtId"=>5293,
#      "currentTId"=>6010,
#      "defaultWin"=>false,
#      "map"=>"de_dust2",
#      "mapOver"=>false}},
#  "listId"=>2337391,
#  "wins"=>{},
#  "liveLog"=>
#   {""=>true,
#    "IrregularTeamKillsRequirement"=>true,
#    "PlayersRequirement"=>true,
#    "NoSuspectEventsInFirstRoundRequirement"=>true,
#    "NotKnifeRoundRequirement"=>true,
#    "BombInPlayRequirement"=>true,
#    "KillsInFirstRoundRequirement"=>true,
#    "FiveKillsWhenEnemyElliminatedRequirement less than 5 five kills in round(s) []"=>true,
#    "MatchStartRequirement"=>true,
#    "MapNameRequirement"=>true,
#    "NoDrawRoundsRequirement"=>true,
#    "FirstRoundOverRequirement"=>true,
#    "RoundOneMaxEquipmentValueRequirement"=>true},
#  "forcedLive"=>false,
#  "forcedDead"=>false}
